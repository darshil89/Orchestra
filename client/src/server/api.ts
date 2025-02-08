"use client"
import { Task } from "@/types/task";
import { API_URL } from "@/utils/db";


export async function postTask(task: Task): Promise<Task> {
  const response = await fetch(`${API_URL}/task`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(task),
  });

  if (!response.ok) {
    throw new Error("Failed to create task");
  }

  const data = await response.json();

  return data;
}