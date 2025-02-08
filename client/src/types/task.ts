export type Task = {
  id: number;
  title: string;
  description: string;
  status: TaskStatus;
  function: {
    url: string;
  }[];
}[];

type TaskStatus = "done" | "pending" | "in-progress";

