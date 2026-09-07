import React, { useState } from "react";
import { 
  X, 
  History, 
  CheckCircle, 
  Clock, 
  ShieldCheck, 
  AlertCircle, 
  Filter,
  FileCheck,
  Send,
  Eye,
  Layers
} from "lucide-react";
import { AuditRecord } from "../types";

interface AuditDrawerProps {
  isOpen: boolean;
  onClose: () => void;
  auditEvents: AuditRecord[];
}

export const AuditDrawer: React.FC<AuditDrawerProps> = ({
  isOpen,
  onClose,
  auditEvents,
}) => {
  const [filterType, setFilterType] = useState<string>("ALL");

  if (!isOpen) return null;

  const filteredEvents = auditEvents.filter((ev) => {
    if (filterType === "ALL") return true;
    return ev.eventType === filterType;
  });

  const getEventBadge = (type: string) => {
    switch (type) {
      case "study_received":
        return {
          icon: <Clock className="w-3.5 h-3.5 text-sky-400" />,
          label: "study_received",
          color: "bg-sky-950/80 text-sky-300 border-sky-800",
        };
      case "image_loaded":
        return {
          icon: <Layers className="w-3.5 h-3.5 text-indigo-400" />,
          label: "image_loaded",
          color: "bg-indigo-950/80 text-indigo-300 border-indigo-800",
        };
      case "image_validated":
        return {
          icon: <ShieldCheck className="w-3.5 h-3.5 text-emerald-400" />,
          label: "image_validated",
          color: "bg-emerald-950/80 text-emerald-300 border-emerald-800",
        };
      case "generic_reading_created":
        return {
          icon: <Eye className="w-3.5 h-3.5 text-amber-400" />,
          label: "generic_reading_created",
          color: "bg-amber-950/80 text-amber-300 border-amber-800",
        };
      case "result_generated":
        return {
          icon: <FileCheck className="w-3.5 h-3.5 text-cyan-400" />,
          label: "result_generated",
          color: "bg-cyan-950/80 text-cyan-300 border-cyan-800",
        };
      case "delivery_started":
      case "delivery_completed":
        return {
          icon: <Send className="w-3.5 h-3.5 text-emerald-400" />,
          label: type,
          color: "bg-emerald-950/80 text-emerald-300 border-emerald-800",
        };
      case "delivery_failed":
        return {
          icon: <AlertCircle className="w-3.5 h-3.5 text-rose-400" />,
          label: "delivery_failed",
          color: "bg-rose-950/80 text-rose-300 border-rose-800",
        };
      default:
        return {
          icon: <CheckCircle className="w-3.5 h-3.5 text-slate-400" />,
          label: type,
          color: "bg-slate-800 text-slate-300 border-slate-700",
        };
    }
  };

  return (
    <div className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex justify-end">
      <div className="bg-slate-900 border-l border-slate-800 w-full max-w-lg h-full p-6 flex flex-col space-y-4 shadow-2xl overflow-hidden animate-in slide-in-from-right duration-200">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-slate-800 pb-3">
          <div className="flex items-center space-x-2">
            <History className="w-5 h-5 text-amber-400" />
            <div>
              <h3 className="text-base font-bold text-white uppercase tracking-tight">
                PISTA DE AUDITORÍA (SECCIÓN 12)
              </h3>
              <p className="text-[11px] text-slate-400">
                Trazabilidad regulatoria • Registro seguro sin PII sensible
              </p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="text-slate-400 hover:text-white p-1 rounded-md"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Filters */}
        <div className="flex items-center space-x-2 overflow-x-auto pb-1 text-xs">
          <Filter className="w-3.5 h-3.5 text-slate-400 flex-shrink-0" />
          {["ALL", "study_received", "image_validated", "generic_reading_created", "result_generated", "delivery_completed"].map((filter) => (
            <button
              key={filter}
              onClick={() => setFilterType(filter)}
              className={`px-2 py-1 rounded font-mono text-[11px] transition whitespace-nowrap ${
                filterType === filter
                  ? "bg-sky-600 text-white font-bold"
                  : "bg-slate-800 text-slate-400 hover:text-slate-200"
              }`}
            >
              {filter === "ALL" ? "Todos los eventos" : filter}
            </button>
          ))}
        </div>

        {/* Events List */}
        <div className="flex-1 overflow-y-auto space-y-2.5 pr-1 font-sans">
          {filteredEvents.length === 0 ? (
            <div className="h-64 flex items-center justify-center text-slate-500 text-xs">
              No hay eventos registrados bajo este filtro.
            </div>
          ) : (
            filteredEvents.map((ev) => {
              const badge = getEventBadge(ev.eventType);
              return (
                <div
                  key={ev.id}
                  className="bg-slate-950 p-3 rounded-lg border border-slate-800/90 text-xs space-y-1.5 hover:border-slate-700 transition"
                >
                  <div className="flex items-center justify-between">
                    <span
                      className={`inline-flex items-center space-x-1 px-2 py-0.5 rounded text-[11px] font-mono border ${badge.color}`}
                    >
                      {badge.icon}
                      <span>{badge.label}</span>
                    </span>
                    <span className="text-[10px] text-slate-500 font-mono">
                      {new Date(ev.timestamp).toLocaleTimeString()}
                    </span>
                  </div>

                  <p className="text-slate-200 text-xs leading-relaxed">
                    {ev.details}
                  </p>

                  <div className="flex justify-between items-center text-[10px] text-slate-500 pt-1 border-t border-slate-900 font-mono">
                    <span>Estudio: {ev.studyId}</span>
                    <span>Actor: {ev.actor}</span>
                  </div>
                </div>
              );
            })
          )}
        </div>

        {/* Footer */}
        <div className="pt-3 border-t border-slate-800 flex justify-between items-center text-[11px] text-slate-400">
          <span>{auditEvents.length} eventos registrados en sesión</span>
          <button
            onClick={onClose}
            className="px-3 py-1.5 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium text-xs"
          >
            Cerrar
          </button>
        </div>
      </div>
    </div>
  );
};
