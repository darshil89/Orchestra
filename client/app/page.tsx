"use client";

import { useState, useEffect } from "react";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Music, Clock, CheckCircle, AlertCircle, Play, RefreshCw } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Task, Tasks, TaskStatus } from "@/types/task";
import { task } from "@/data/task";
import { postTask } from "@/server/api";

const initialTasks: Tasks = task;

export default function Home() {
  const [logs, setLogs] = useState<string[]>([]);

  const [tasks, setTasks] = useState<Tasks>(
    initialTasks.map(task => ({ ...task, status: "stale" }))
  );

  const startTasks = async () => {

    const response = await postTask(tasks)

    console.log("response from server:", response)
  };

  const resetTasks = () => {
    setTasks(initialTasks.map(task => ({ ...task, status: "stale" })));
    setLogs(prev => [
      ...prev,
      `[${new Date().toISOString().replace('T', ' ').slice(0, 19)}] INFO: Tasks reset to initial state`
    ]);
  };

  const getStatusColor = (status: TaskStatus) => {
    switch (status) {
      case "done":
        return "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300";
      case "in-progress":
        return "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300";
      case "stale":
        return "bg-amber-100 text-amber-800 dark:bg-amber-900 dark:text-amber-300";
      default:
        return "bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300";
    }
  };

  const getStatusIcon = (status: TaskStatus) => {
    switch (status) {
      case "done":
        return <CheckCircle className="h-4 w-4 text-green-600 dark:text-green-400" />;
      case "in-progress":
        return <Clock className="h-4 w-4 text-blue-600 dark:text-blue-400" />;
      case "stale":
        return <AlertCircle className="h-4 w-4 text-amber-600 dark:text-amber-400" />;
      default:
        return null;
    }
  };

  const getAnimationClass = (task: Task) => {
    if (task.animationState === "running") {
      return "animate-pulse-wave";
    } else if (task.animationState === "completed" && task.status === "done") {
      return "animate-success-glow";
    }
    return "";
  };

  const getProgressBarClass = (task: Task) => {
    if (task.status === "done") {
      return "bg-gradient-to-r from-teal-400 to-green-500 w-full";
    } else if (task.status === "in-progress") {
      return task.animationState === "running"
        ? "bg-gradient-to-r from-blue-400 to-cyan-300 animate-progress-wave w-3/3"
        : "bg-blue-500 w-3/3";
    } else {
      return "bg-amber-500 w-3/3";
    }
  };

  useEffect(() => {
    // connect the websocket 

    const ws = new WebSocket("ws://localhost:8080/ws");

    ws.onopen = () => {
      console.log("Connected to websocket");
    }

    ws.onmessage = (event) => {
      console.log("Message from server:", event.data);

      // update the logs 
      setLogs((prev) => [
        ...prev,
        `[${new Date().toISOString().replace('T', ' ').slice(0, 19)}] ${event.data}`
      ]);
      const data = JSON.parse(event.data);

      setTasks((prev) =>
        prev.map((task) => {
          if (task.id === data.id) {
            // set animation here 
            if (data.status === "in-progress") {
              return { ...task, status: data.status, animationState: "running" }
            }
            if (data.status === "done") {
              return { ...task, status: data.status, animationState: "completed" }
            }
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
    <div className="min-h-screen bg-gradient-to-b from-slate-900 to-slate-800 p-6 md:p-10">
      <div className="max-w-7xl mx-auto">
        <header className="mb-12">
          <h1 className="text-4xl md:text-5xl font-bold tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-cyan-300 via-blue-500 to-teal-400 animate-text-shimmer">
            Orchestra — Distributed Job Orchestration
          </h1>
          <p className="mt-4 text-cyan-100 text-lg">
            Manage and monitor your distributed tasks with precision
          </p>
        </header>

        <div className="flex flex-col md:flex-row gap-4 mb-8">
          <Button
            onClick={startTasks}
            className="bg-gradient-to-r from-cyan-500 to-blue-600 hover:from-cyan-600 hover:to-blue-700 text-white shadow-lg shadow-cyan-500/20"
          >
            <Play className="mr-2 h-4 w-4" /> Start
          </Button>
          <Button
            onClick={resetTasks}
            variant="outline"
            className="border-cyan-500 text-white bg-cyan-900"
          >
            <RefreshCw className="mr-2 h-4 w-4" /> Reset Tasks
          </Button>
        </div>

        <section className="mb-10">
          <h2 className="text-2xl font-semibold mb-6 text-cyan-100">Active Tasks</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {tasks.map((task) => (
              <TooltipProvider key={task.id} delayDuration={300}>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <div
                      className={`bg-slate-800/80 backdrop-blur-sm rounded-lg p-6 shadow-md hover:shadow-lg transition-all duration-200 border border-cyan-900/50 ${getAnimationClass(task)}`}
                    >
                      <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center space-x-3">
                          <h3 className="text-xl font-medium text-cyan-50">{task.title}</h3>
                        </div>
                        <Badge className={`flex items-center gap-1.5 ${getStatusColor(task.status)}`}>
                          {getStatusIcon(task.status)}
                          {task.status}
                        </Badge>
                      </div>
                      <div className="h-2 bg-slate-700 rounded-full overflow-hidden">
                        <div
                          className={`h-full ${getProgressBarClass(task)}`}
                        />
                      </div>
                    </div>
                  </TooltipTrigger>
                  <TooltipContent side="bottom" className="p-3 max-w-xs bg-slate-800 border-cyan-900/50">
                    <p className="text-cyan-50">{task.description}</p>
                  </TooltipContent>
                </Tooltip>
              </TooltipProvider>
            ))}
          </div>
        </section>

        <section>
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-2xl font-semibold text-cyan-100">System Logs</h2>
            <Badge variant="outline" className="text-xs border-cyan-700 text-cyan-300">
              Live Updates
            </Badge>
          </div>
          <div className="bg-slate-950 rounded-lg border border-cyan-900/30 shadow-xl">
            <div className="flex items-center px-4 py-2 border-b border-cyan-900/30 bg-slate-900">
              <div className="flex space-x-2">
                <div className="w-3 h-3 rounded-full bg-red-500"></div>
                <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
                <div className="w-3 h-3 rounded-full bg-green-500"></div>
              </div>
              <div className="ml-4 text-sm text-cyan-400">orchestra-logs</div>
            </div>
            <ScrollArea className="h-[350px] p-4 font-mono text-sm">
              {logs.map((log, index) => {
                const isError = log.includes("ERROR");
                const isWarning = log.includes("WARNING");
                const isSuccess = log.includes("SUCCESS");

                let logClass = "text-cyan-300";
                if (isError) logClass = "text-red-400";
                else if (isWarning) logClass = "text-yellow-400";
                else if (isSuccess) logClass = "text-emerald-400";

                return (
                  <div key={index} className={`py-1 ${logClass}`}>
                    {log}
                  </div>
                );
              })}
            </ScrollArea>
          </div>
        </section>
      </div>
    </div>
  );
}