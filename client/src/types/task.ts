export type Task = {
  id: number;
  title: string;
  description: string;
  status: TaskStatus;
  function: {
    url: string;
  }[];
}

export type Tasks = Task[];

type TaskStatus = "done" | "stale" | "in-progress";

