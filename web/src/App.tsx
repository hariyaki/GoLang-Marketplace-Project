import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'

interface Listing {
  id: number
  title: string
  description: string
  price_jpy: number
  image_url?: string
  created_at: string
}

const API = import.meta.env.VITE_API as string || 'http://localhost:8080'

export default function App() {
  const [listings, setListings] = useState<Listing[]>([])
  const [form, setForm] = useState({ title: '', description: '', price_jpy: '' })
  const [msg, setMsg] = useState('')

  const fetchListings = async () => {
    const res = await fetch(`${API}/listings`)
    setListings(await res.json())
  }

  useEffect(() => { fetchListings() }, [])

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

  return (
    <div style={{ maxWidth: 600, margin: '2rem auto', fontFamily: 'sans-serif' }}>
      <h1>Go Marketplace</h1>

      <form onSubmit={submit}>
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
      <ul>
        {listings.map(l => (
          <li key={l.id}>
            <strong>{l.title}</strong> – ¥{l.price_jpy}<br />
            {l.description}
          </li>
        ))}
      </ul>
    </div>
  )
}
