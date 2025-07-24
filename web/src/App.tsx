import { useEffect, useState } from 'react'
import type { FormEvent, ChangeEvent } from 'react'
import { API } from './lib/api'

interface Listing {
  id: number
  title: string
  description: string
  price_jpy: number
  image_url?: string
  created_at: string
}

export default function App() {
  const [listings, setListings] = useState<Listing[]>([])
  const [form, setForm] = useState({ title: '', description: '', price_jpy: '' })
  const [msg, setMsg] = useState('')

  // ---------------- fetch helpers ----------------
  const fetchListings = async () => {
    const res = await fetch(`${API}/listings`)
    setListings(await res.json())
  }

  const uploadImage = async (id: number, file: File) => {
    const fd = new FormData()
    fd.append('file', file)
    await fetch(`${API}/listings/${id}/image`, { method: 'PUT', body: fd })
    fetchListings()
  }

  useEffect(() => { fetchListings() }, [])

  // ---------------- form submit ----------------
  const submit = async (e: FormEvent) => {
    e.preventDefault()
    const res = await fetch(`${API}/listings`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ ...form, price_jpy: Number(form.price_jpy) })
    })
    if (res.ok) {
      setMsg('Created!')
      setForm({ title: '', description: '', price_jpy: '' })
      fetchListings()
    } else {
      setMsg('Error creating listing')
    }
  }

  // ---------------- render ----------------
  return (
    <div style={{ maxWidth: 700, margin: '2rem auto', fontFamily: 'sans-serif' }}>
      <h1>Go Marketplace</h1>

      <form onSubmit={submit} style={{ marginBottom: '1rem' }}>
        <input
          placeholder="Title"
          required
          value={form.title}
          onChange={e => setForm({ ...form, title: e.target.value })}
        />
        <br />
        <textarea
          placeholder="Description"
          required
          value={form.description}
          onChange={e => setForm({ ...form, description: e.target.value })}
        />
        <br />
        <input
          type="number"
          placeholder="Price (JPY)"
          required
          value={form.price_jpy}
          onChange={e => setForm({ ...form, price_jpy: e.target.value })}
        />
        <button type="submit">Create</button>
      </form>

      {msg && <p>{msg}</p>}

      <h2>Listings</h2>
      <ul style={{ listStyle: 'none', padding: 0 }}>
        {listings.map(l => (
          <li key={l.id} style={{ marginBottom: '1.5rem' }}>
            {l.image_url && (
              <img
                src={`${API}${l.image_url}`}
                alt={l.title}
                width={150}
                style={{ display: 'block', marginBottom: '0.5rem' }}
              />
            )}

            <strong>{l.title}</strong> – ¥{l.price_jpy}<br />
            {l.description}
            <br />

            <input
              type="file"
              accept="image/*"
              onChange={(e: ChangeEvent<HTMLInputElement>) => {
                const f = e.target.files?.[0]
                if (f) uploadImage(l.id, f)
              }}
            />
          </li>
        ))}
      </ul>
    </div>
  )
}
