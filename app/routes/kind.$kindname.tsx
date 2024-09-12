import { useParams } from '@remix-run/react';
import useSWR from 'swr';
import semver from 'semver';

type GVK = {
  group: string;
  version: string;
  kind: string;
};

type GVKAndOccurrence = {
  GVK: GVK;
  Lifecycles: { kubernetesMinorRelease: string; APILifecycle: string }[];
};

type GVKSpan = {
  GVK: GVK;
  spans: { start: string; end: string; lifecycle: string }[];
};

const fetcher = (url: string) => fetch(url).then(res => res.json());

export default function KindDetail() {
  const { kindname } = useParams();
  const { data, error } = useSWR<GVKAndOccurrence[]>(`/api/kind/${kindname}`, fetcher);

  if (error) return <div>Failed to load</div>;
  if (!data) return <div>Loading...</div>;

  const gvkStrings = data.map(gvk => `${gvk.GVK.group}/${gvk.GVK.version} ${gvk.GVK.kind}`);

  const versions = Array.from(new Set(data.flatMap(gvk => gvk.Lifecycles.map(l => l.kubernetesMinorRelease))))
    .sort((a, b) => {
      const semverA = semver.coerce(a);
      const semverB = semver.coerce(b);
      return semverA!.compare(semverB!);
    });

  const spans: GVKSpan[] = data.map(gvk => {
    const spans = [] as { start: string; end: string; lifecycle: string }[];
    let start = '';
    let lifecycle = '';
    gvk.Lifecycles.forEach((item, index, arr) => {
      if (index === 0) {
        start = item.kubernetesMinorRelease;
        lifecycle = item.APILifecycle;
      }
      if (lifecycle !== item.APILifecycle || index === arr.length - 1) {
        spans.push({ start, end: item.kubernetesMinorRelease, lifecycle });
        start = item.kubernetesMinorRelease;
        lifecycle = item.APILifecycle;
      }
    });
    return { GVK: gvk.GVK, spans };
  });

  return (
    <div className="p-4">
      <h1 className="text-3xl mb-4">
        <span className="text-gray-500">Availability for </span>
        <span>{kindname}</span>
      </h1>
      <div className="flex">
        <div className="flex-none text-right p-4">
          {gvkStrings.map((item, index) => (
            <div key={index} className="p-2 h-10">
              <span>{item}</span>
            </div>
          ))}
          <div className="p-2 h-10">
            <a className="text-blue-500 hover:underline" target="_blank" href="https://endoflife.date/kubernetes">
              Kubernetes Version
            </a>
          </div>
        </div>
        <div className="flex-1 p-4">
          {spans.map((span, spanIndex) => (
            <div
              key={spanIndex}
              className="grid h-10"
              style={{
                gridTemplateColumns: `repeat(${versions.length}, minmax(0, 1fr))`,
              }}
            >
              {span.spans.map((item, itemIndex) => (
                <div
                  key={itemIndex}
                  className={`flex border-neutral ${
                    itemIndex !== span.spans.length - 1 ? 'border-l' : 'border-x'
                  }`}
                  style={{
                    gridColumn: `${versions.indexOf(item.start) + 1} / ${versions.indexOf(item.end) + 2}`,
                  }}
                >
                  <div
                    className={`w-full h-4 my-auto ${
                      item.lifecycle === 'stable' ? 'bg-green-500' : 'bg-yellow-500'
                    }`}
                  ></div>
                </div>
              ))}
            </div>
          ))}
          <div
            className="h-10 grid text-center divide-x divide-neutral"
            style={{
              gridTemplateColumns: `repeat(${versions.length}, minmax(0, 1fr))`,
            }}
          >
            {versions.map((item, index) => (
              <div key={index} className="p-2">
                <span>{item}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
