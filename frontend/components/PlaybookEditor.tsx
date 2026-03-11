"use client";

import { useState } from "react";

interface Node {
  id: string;
  type: "trigger" | "condition" | "action" | "delay";
  config: Record<string, any>;
  next: string[];
  x: number;
  y: number;
}

interface Playbook {
  id: string;
  name: string;
  description: string;
  startNode: string;
  nodes: Node[];
}

export default function PlaybookEditor() {
  const [playbook, setPlaybook] = useState<Playbook>({
    id: "demo",
    name: "Demo Playbook",
    description: "Drag-and-drop playbook",
    startNode: "trigger1",
    nodes: [
      { id: "trigger1", type: "trigger", config: {}, next: ["action1"], x: 50, y: 50 },
      { id: "action1", type: "action", config: { action: "isolate_host" }, next: [], x: 250, y: 50 },
    ],
  });

  const addNode = (type: Node["type"]) => {
    const newNode: Node = {
      id: `${type}${Date.now()}`,
      type,
      config: type === "delay" ? { seconds: 30 } : type === "action" ? { action: "block_ip" } : {},
      next: [],
      x: 150 + Math.random() * 200,
      y: 150 + Math.random() * 200,
    };
    setPlaybook({ ...playbook, nodes: [...playbook.nodes, newNode] });
  };

  const savePlaybook = async () => {
    await fetch("/api/v1/playbooks", {
      method: "POST",
      headers: { "Content-Type": "application/json", "Authorization": `Bearer ${localStorage.getItem("token")}` },
      body: JSON.stringify(playbook),
    });
    alert("Playbook saved");
  };

  return (
    <main className="p-8">
      <h1 className="text-3xl font-bold mb-6">Playbook Editor</h1>

      <section className="mb-4 flex gap-2">
        <button onClick={() => addNode("trigger")} className="px-3 py-2 bg-gray-600 text-white rounded">+ Trigger</button>
        <button onClick={() => addNode("condition")} className="px-3 py-2 bg-blue-600 text-white rounded">+ Condition</button>
        <button onClick={() => addNode("action")} className="px-3 py-2 bg-green-600 text-white rounded">+ Action</button>
        <button onClick={() => addNode("delay")} className="px-3 py-2 bg-yellow-600 text-white rounded">+ Delay</button>
        <button onClick={savePlaybook} className="px-4 py-2 bg-phoenix-orange text-white rounded ml-auto">Save</button>
      </section>

      <section className="relative bg-gray-50 border rounded" style={{ height: "500px" }}>
        <svg className="absolute inset-0 w-full h-full">
          {playbook.nodes.map((node) =>
            node.next.map((targetId) => {
              const target = playbook.nodes.find((n) => n.id === targetId);
              if (!target) return null;
              return (
                <line
                  key={`${node.id}-${targetId}`}
                  x1={node.x + 60}
                  y1={node.y + 30}
                  x2={target.x}
                  y2={target.y + 30}
                  stroke="#888"
                  strokeWidth="2"
                  markerEnd="url(#arrowhead)"
                />
              );
            })
          )}
          <defs>
            <marker id="arrowhead" markerWidth="10" markerHeight="7" refX="9" refY="3.5" orient="auto">
              <polygon points="0 0, 10 3.5, 0 7" fill="#888" />
            </marker>
          </defs>
        </svg>
        {playbook.nodes.map((node) => (
          <div
            key={node.id}
            className="absolute bg-white border rounded p-2 shadow cursor-move"
            style={{ left: node.x, top: node.y, width: "120px" }}
            draggable
            onDragEnd={(e) => {
              const updated = playbook.nodes.map((n) =>
                n.id === node.id ? { ...n, x: e.clientX - 60, y: e.clientY - 30 } : n
              );
              setPlaybook({ ...playbook, nodes: updated });
            }}
          >
            <div className="font-semibold text-sm capitalize">{node.type}</div>
            {node.type === "action" && <div className="text-xs text-gray-600">{node.config.action}</div>}
            {node.type === "delay" && <div className="text-xs text-gray-600">{node.config.seconds}s</div>}
          </div>
        ))}
      </section>
    </main>
  );
}
