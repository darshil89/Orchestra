// pushing random data to github.
// pushing an image to dockerhub.
// calling an api.

import { Task } from "@/types/task";

export const task: Task = [
  {
    id: 1,
    title: "DB Push",
    description: "Push a random data to the database.",
    status: "stale",
    function: [
      {
        url: "",
      },
    ],
  },{
    id: 2,
    title: "Dockerhub Push",
    description: "Push an public image to Dockerhub.",
    status: "stale",
    function: [
      {
        url: "",
        
      },
      {
        url: "mahraurdarshil89/testing",
      }
    ],
  },{
    id: 3,
    title: "API Call",
    description: "Make a GET/POST request to an API.",
    status: "stale",
    function: [
      {
        url: "",
      },
    ],
  }
];
