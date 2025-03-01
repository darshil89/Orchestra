import React from "react";
import { cn } from "@/lib/utils";

interface TerminalProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
}

export function Terminal({ children, className, ...props }: TerminalProps) {
  return (
    <div
      className={cn(
        "rounded-lg border bg-card text-card-foreground shadow-sm",
        className
      )}
      {...props}
    >
      <div className="flex items-center border-b bg-muted/50 px-3 py-2">
        <div className="flex h-2 w-2 rounded-full bg-red-500" />
        <div className="ml-2 flex h-2 w-2 rounded-full bg-yellow-500" />
        <div className="ml-2 flex h-2 w-2 rounded-full bg-green-500" />
        <div className="ml-4 text-xs text-muted-foreground">Terminal</div>
      </div>
      <div className="p-4 font-mono text-sm">{children}</div>
    </div>
  );
}