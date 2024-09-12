import { useParams } from '@remix-run/react'
import semver from 'semver'
import useSWR from 'swr'

type FieldSummary = {
  FieldPath: string
  LifeCycles: {
    KubernetesMinorRelease: string
    APILifecycle: string
  }[]
}

type ApiResponse = {
  Summary: Record<string, FieldSummary>
}

const fetcher = (url: string) => fetch(url).then((res) => res.json())

export default function FieldGVK() {
  const { gvk } = useParams()
  const { data, error } = useSWR<ApiResponse>(`/api/field/${encodeURIComponent(gvk as string)}`, fetcher)

  if (error) return <div>Failed to load</div>
  if (!data) return <div>Loading...</div>

  const fieldSummary: FieldSummary[] = Object.values(data.Summary)

  const versions = Array.from(
    new Set(fieldSummary.flatMap((field) => field.LifeCycles.map((lc) => lc.KubernetesMinorRelease)))
  ).sort((a, b) => {
    const semverA = semver.coerce(a)
    const semverB = semver.coerce(b)
    return semverB!.compare(semverA!)
  })

  const aggregatedFields = new Map<string, Map<string, string>>()
  fieldSummary.forEach((field) => {
    const fieldPath = field.FieldPath
    const fieldMap = new Map<string, string>()
    field.LifeCycles.forEach((lc) => {
      fieldMap.set(lc.KubernetesMinorRelease, lc.APILifecycle)
    })
    aggregatedFields.set(fieldPath, fieldMap)
  })

  const fields = Array.from(aggregatedFields.keys()).sort()

  return (
    <div>
      <h1 className="text-3xl p-4">
        <span className="text-gray-500">Fields Availability for </span>
        <span>{gvk}</span>
      </h1>
      <div className="overflow-x-auto overflow-y-auto border-neutral" style={{ height: '85vh' }}>
        <table className="table table-compact w-full">
          <thead>
            <tr>
              <th className="border px-4 py-2">Field</th>
              {versions.map((version) => (
                <th key={version} className="border px-4 py-2">
                  {version}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {fields.map((field) => (
              <tr key={field}>
                <td className="border px-4 py-2">
                  {field.replace('io.k8s.api.', '').replace((gvk as string).split(' ')[0].replace('/', '.') + '.', '')}
                </td>
                {versions.map((version) => {
                  const lifecycle = aggregatedFields.get(field)?.get(version)
                  return (
                    <td
                      key={version}
                      className={`border px-4 py-2 ${
                        lifecycle === 'stable'
                          ? 'bg-green-500'
                          : lifecycle === 'deprecated'
                          ? 'bg-yellow-500'
                          : 'bg-neutral'
                      }`}
                    >
                      {lifecycle}
                    </td>
                  )
                })}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
