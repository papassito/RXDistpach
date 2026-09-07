import React from "react";
import { 
  CheckCircle2, 
  Circle, 
  Layers, 
  Eye, 
  Package, 
  Send, 
  FileCheck 
} from "lucide-react";
import { StudyStatus } from "../types";

interface WorkflowPipelineProps {
  currentStatus: StudyStatus;
  onSelectStage: (stageIndex: number) => void;
  activeStage: number;
}

export const WorkflowPipeline: React.FC<WorkflowPipelineProps> = ({
  currentStatus,
  onSelectStage,
  activeStage,
}) => {
  const steps = [
    {
      id: 0,
      title: "1. Recepción",
      subtitle: "Study Ingest",
      icon: <Layers className="w-4 h-4" />,
      isCompleted: true,
    },
    {
      id: 1,
      title: "2. Image Engine",
      subtitle: "Integridad SHA-256",
      icon: <Layers className="w-4 h-4" />,
      isCompleted: currentStatus !== "RECEIVED",
    },
    {
      id: 2,
      title: "3. Generic Reader",
      subtitle: "Lectura Genérica",
      icon: <Eye className="w-4 h-4" />,
      isCompleted: currentStatus === "GENERIC_READ_COMPLETED" || currentStatus === "PACKAGE_BUILT" || currentStatus === "DELIVERED",
    },
    {
      id: 3,
      title: "4. Result Builder",
      subtitle: "Paquete RX_RESULT",
      icon: <Package className="w-4 h-4" />,
      isCompleted: currentStatus === "PACKAGE_BUILT" || currentStatus === "DELIVERED",
    },
    {
      id: 4,
      title: "5. Delivery",
      subtitle: "Despacho y Envío",
      icon: <Send className="w-4 h-4" />,
      isCompleted: currentStatus === "DELIVERED",
    },
  ];

  return (
    <div className="bg-slate-900 border border-slate-800 rounded-xl p-3 shadow-md">
      <div className="flex items-center justify-between overflow-x-auto gap-2 py-1">
        {steps.map((step, idx) => {
          const isActive = activeStage === step.id;
          return (
            <React.Fragment key={step.id}>
              <button
                onClick={() => onSelectStage(step.id)}
                className={`flex items-center space-x-2.5 px-3 py-2 rounded-lg transition text-left flex-shrink-0 ${
                  isActive
                    ? "bg-sky-950/80 border border-sky-500/70 text-white shadow-sm"
                    : "hover:bg-slate-800/60 text-slate-400"
                }`}
              >
                <div
                  className={`w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold ${
                    step.isCompleted
                      ? "bg-emerald-500/20 text-emerald-400 border border-emerald-500/40"
                      : isActive
                      ? "bg-sky-500/20 text-sky-400 border border-sky-500/40"
                      : "bg-slate-800 text-slate-500 border border-slate-700"
                  }`}
                >
                  {step.isCompleted ? (
                    <CheckCircle2 className="w-3.5 h-3.5" />
                  ) : (
                    <span>{step.id + 1}</span>
                  )}
                </div>
                <div>
                  <div className={`text-xs font-semibold ${isActive ? "text-sky-200 font-bold" : "text-slate-200"}`}>
                    {step.title}
                  </div>
                  <div className="text-[10px] text-slate-500 truncate max-w-[110px]">
                    {step.subtitle}
                  </div>
                </div>
              </button>

              {idx < steps.length - 1 && (
                <div className="h-0.5 w-4 bg-slate-800 hidden sm:block flex-shrink-0" />
              )}
            </React.Fragment>
          );
        })}
      </div>
    </div>
  );
};
