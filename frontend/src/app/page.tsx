"use client";
import Link from "next/link";

export default function Home() {
  return (
    <div className="flex flex-col items-center justify-center h-screen space-y-4">
      <h1 className="text-4xl font-bold">Welcome to AI Task Manager</h1>
      <Link href="/login" className="text-blue-500 underline">
        Go to Login
      </Link>
      <Link href="/signup" className="text-blue-500 underline">
        Go to Signup
      </Link>
      <Link href="/tasks" className="text-blue-500 underline">
        Go to Tasks
      </Link>
    </div>
  );
}
