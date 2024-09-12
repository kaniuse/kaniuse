import { useState } from 'react';
import { Link } from '@remix-run/react';
import useSWR from 'swr';

const fetcher = () => fetch('/api/kinds').then(res => res.json());

export default function Index() {
  const { data: kinds, error } = useSWR<string[]>('kinds', fetcher);

  const [filter, setFilter] = useState('');

  const filteredKinds = kinds?.filter((kind) =>
    kind.toLowerCase().includes(filter.toLowerCase())
  ) || [];

  return (
    <div className="flex justify-left items-center m-4">
      <div>
        <div className="p-4">
          <span className="text-lg">Can I use</span>
          <input
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            placeholder="Kind, eg. Ingress"
            className="input input-bordered mx-5"
          />
          <span className="text-lg">?</span>
        </div>
        <div className="p-4">
          <ul>
            {filteredKinds.map((kind, index) => (
              <li key={index}>
                <Link to={`/kind/${kind}`} target="_blank">
                  {kind}
                </Link>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </div>
  );
}
