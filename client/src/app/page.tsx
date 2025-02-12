"use client";
import { task } from "@/data/task";
import { postTask } from "@/server/api";
import { Tasks } from "@/types/task";
import { useState, useEffect } from "react";

const initialTasks: Tasks = task;

const statusColors: Record<string, string> = {
  stale: "border-gray-400",
  "in-progress": "border-yellow-400",
  done: "border-green-400",
};

export default function OrchestraLanding() {
  const [tasks, setTasks] = useState<Tasks>(
    initialTasks.map(task => ({ ...task, status: "stale" }))
  );

  const startTasks = async () => {

    const response = await postTask(tasks)

    console.log("response from server:", response)
  };

  const resetTasks = () => {
    setTasks(initialTasks.map(task => ({ ...task, status: "stale" })));
  };

  useEffect(() => {
    // connect the websocket 

    const ws = new WebSocket("ws://localhost:8080/ws");

    ws.onopen = () => {
      console.log("Connected to websocket");
    }

    ws.onmessage = (event) => {
      console.log("Message from server:", event.data);
      const data = JSON.parse(event.data);

      setTasks((prev) =>
        prev.map((task) => {
          if (task.id === data.id) {
            return { ...task, status: data.status }
          }
          return task

        }
        )
      )

    }

    ws.onerror = (error) => {
      console.error("⚠️ WebSocket Error:", error);
    };

    ws.onclose = () => {
      console.warn("⚠️ WebSocket Disconnected, attempting to reconnect...");

    };

    return () => {
      ws.close();
    };

  }, []);

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-900 text-white p-6">
      <h1 className="text-4xl font-bold mb-6">Orchestra - Distributed Job Orchestration</h1>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {tasks.map((task) => (
          <div
            key={task.id}
            className={`relative p-6 border-4 ${statusColors[task.status]} rounded-xl transition-all duration-500 hover:shadow-lg hover:scale-105 group`}
          >
            <h2 className="text-2xl font-semibold">{task.title}</h2>
            <p className="text-lg text-gray-300">Status: {task.status}</p>
            <div className="absolute inset-0 rounded-xl bg-black p-6 flex flex-col justify-center opacity-0 group-hover:opacity-100 transition-all duration-300">
              <p>{task.description}</p>
              <p className="mt-2 text-sm text-gray-400">URLs: {task.function.length}</p>
            </div>
          </div>
        ))}
      </div>
      <div className="mt-8 flex gap-4">
        <button
          onClick={startTasks}
          className="px-6 py-3 bg-blue-500 hover:bg-blue-600 text-white font-bold rounded-xl transition-all"
        >
          Start
        </button>
        <button
          onClick={resetTasks}
          className="px-6 py-3 bg-red-500 hover:bg-red-600 text-white font-bold rounded-xl transition-all"
        >
          Reset
        </button>
      </div>
    </div>
  );
}
