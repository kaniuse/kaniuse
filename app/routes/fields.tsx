import { Link } from '@remix-run/react'
import { useEffect, useState } from 'react'

export default function Fields() {
  const [filter, setFilter] = useState('')
  const [gvks, setGvks] = useState<string[]>([])

  useEffect(() => {
    const fetchGvks = async () => {
      const response = await fetch('/api/gvks')
      const data = await response.json()
      setGvks(data)
    }
    fetchGvks()
  }, [])

  const filteredGVKs = filter ? gvks.filter((gvk) => gvk.toLowerCase().includes(filter.toLowerCase())) : gvks

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
            {filteredGVKs.map((gvk, index) => (
              <li key={index}>
                <Link to={`/field/${encodeURIComponent(gvk)}`} target="_blank">
                  {gvk}
                </Link>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </div>
  )
}
