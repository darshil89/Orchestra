// pushing random data to github.
// pushing an image to dockerhub.
// calling an api.

import { Task } from "@/types/task";

const task: Task = [
  {
    id: 1,
    title: "Github Push",
    description: "Push a random data to a public repository on Github.",
    status: "pending",
    function: [
      {
        url: "",
      },
    ],
  },{
    id: 2,
    title: "Dockerhub Push",
    description: "Push an public image to Dockerhub.",
    status: "pending",
    function: [
      {
        url: "",
        
      },
      {
        url: "",
      }
    ],
  },{
    id: 3,
    title: "API Call",
    description: "Make a GET/POST request to an API.",
    status: "pending",
    function: [
      {
        url: "",
      },
    ],
  }
];
