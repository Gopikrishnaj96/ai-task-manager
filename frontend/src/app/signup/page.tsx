"use client";
import { useState } from 'react';
import { useRouter } from 'next/navigation';

export default function Signup() {
  const router = useRouter();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  async function handleSignup(e: React.FormEvent) {
    e.preventDefault();
    const res = await fetch('http://localhost:4000/signup', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password }),
    });
    if (res.ok) {
      router.push('/login');
    } else {
      const data = await res.json();
      alert(data.error || 'Signup failed');
    }
  }

  return (
    <div className="flex items-center justify-center h-screen bg-gray-100">
      <form onSubmit={handleSignup} className="bg-white p-6 rounded shadow-md">
        <h1 className="text-2xl mb-4">Sign Up</h1>
        <input
          type="text"
          placeholder="Username"
          className="border p-2 w-full mb-4"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
        />
        <input
          type="password"
          placeholder="Password"
          className="border p-2 w-full mb-4"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        <button type="submit" className="bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600">
          Sign Up
        </button>
      </form>
    </div>
  );
}
