export type Task = {
    id: number;
    title: string;
    description: string;
    status: TaskStatus;
    animationState: "completed" | "running" | "idle";
    function: {
      url: string;
    }[];
  }
  
  export type Tasks = Task[];
  
 export type TaskStatus = "done" | "stale" | "in-progress";