"use client";
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';

type Task = {
  id: number;
  title: string;
  description: string;
  status: string;
};

export default function Tasks() {
  const router = useRouter();
  const [tasks, setTasks] = useState<Task[]>([]);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) {
      router.push('/login');
      return;
    }
    fetchTasks(token);
  }, []);

  async function fetchTasks(token: string) {
    const res = await fetch('http://localhost:4000/tasks', {
      headers: { 'Authorization': token },
    });
    if (res.ok) {
      const data = await res.json();
      setTasks(data);
    }
  }

  async function handleCreateTask(e: React.FormEvent) {
    e.preventDefault();
    const token = localStorage.getItem('token');
    if (!token) return;
    const res = await fetch('http://localhost:4000/tasks', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': token,
      },
      body: JSON.stringify({ title, description }),
    });
    if (res.ok) {
      setTitle('');
      setDescription('');
      fetchTasks(token);
    }
  }

  function handleLogout() {
    localStorage.removeItem('token');
    router.push('/login');
  }

  return (
    <div className="min-h-screen bg-gray-100 p-6">
      <div className="max-w-xl mx-auto">
        <div className="flex justify-between mb-4">
          <h1 className="text-2xl font-bold">Tasks</h1>
          <button onClick={handleLogout} className="bg-red-500 text-white px-4 py-2 rounded">
            Logout
          </button>
        </div>
        <form onSubmit={handleCreateTask} className="mb-6">
          <input
            type="text"
            placeholder="Title"
            className="border p-2 w-full mb-2"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
          />
          <textarea
            placeholder="Description"
            className="border p-2 w-full mb-2"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
          />
          <button type="submit" className="bg-blue-600 text-white py-2 px-4 rounded hover:bg-blue-700">
            Create Task
          </button>
        </form>
        <ul className="space-y-2">
          {tasks.map(task => (
            <li key={task.id} className="bg-white p-4 rounded shadow">
              <h2 className="font-bold">{task.title}</h2>
              <p>{task.description}</p>
              <p>Status: {task.status}</p>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
