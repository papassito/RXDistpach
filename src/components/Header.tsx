import React from "react";
import { 
  Activity, 
  Smartphone, 
  Monitor, 
  Code, 
  PlusCircle, 
  ShieldCheck,
  History,
  Network
} from "lucide-react";

interface HeaderProps {
  isMobileView: boolean;
  onToggleMobileView: () => void;
  onOpenNewStudy: () => void;
  onOpenAudit: () => void;
  onOpenCodebase: () => void;
  onOpenMesh: () => void;
  auditCount: number;
  meshOnlineCount?: number;
}

export const Header: React.FC<HeaderProps> = ({
  isMobileView,
  onToggleMobileView,
  onOpenNewStudy,
  onOpenAudit,
  onOpenCodebase,
  onOpenMesh,
  auditCount,
  meshOnlineCount = 9,
}) => {
  return (
    <header className="bg-slate-900 border-b border-slate-800 text-white sticky top-0 z-30 shadow-md">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        {/* Brand & Subtitle */}
        <div className="flex items-center space-x-3">
          <div className="w-10 h-10 rounded-lg bg-sky-600 flex items-center justify-center shadow-lg shadow-sky-500/20 ring-1 ring-sky-400/30">
            <Activity className="w-6 h-6 text-white" />
          </div>
          <div>
            <div className="flex items-center space-x-2">
              <span className="font-black text-lg tracking-tight text-slate-100">
                RX DISPATCH
              </span>
              <span className="text-xs font-semibold px-2 py-0.5 rounded bg-sky-500/20 text-sky-300 border border-sky-500/30">
                BY KLIK SOFT PRO
              </span>
              <span className="hidden md:inline-flex items-center space-x-1 text-[11px] font-mono px-2 py-0.5 rounded bg-emerald-950/60 text-emerald-400 border border-emerald-800/40">
                <ShieldCheck className="w-3 h-3" />
                <span>100% Go • Cero Monolito</span>
              </span>
            </div>
            <p className="text-xs text-slate-400 hidden sm:block">
              Sistema de envío de resultados e imágenes de Rayos X (9 Microservicios)
            </p>
          </div>
        </div>

        {/* Action Controls */}
        <div className="flex items-center space-x-2 sm:space-x-3">
          {/* Go Mesh Monitoring */}
          <button
            id="btn-open-mesh"
            onClick={onOpenMesh}
            className="flex items-center space-x-1.5 px-2.5 py-1.5 rounded-md text-xs font-medium bg-emerald-950/70 hover:bg-emerald-900/80 border border-emerald-700/60 transition text-emerald-300 shadow-sm"
            title="Monitorear Malla de 9 Microservicios Go"
          >
            <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
            <Network className="w-3.5 h-3.5" />
            <span className="hidden sm:inline">Malla Go (9/9)</span>
          </button>

          {/* New Study Intake */}
          <button
            id="btn-new-study"
            onClick={onOpenNewStudy}
            className="flex items-center space-x-1.5 px-3 py-1.5 rounded-md text-xs font-medium bg-sky-600 hover:bg-sky-500 active:bg-sky-700 transition shadow-sm"
            title="Ingresar nuevo estudio radiográfico"
          >
            <PlusCircle className="w-4 h-4" />
            <span className="hidden sm:inline">Nuevo Estudio</span>
          </button>

          {/* Audit Logs */}
          <button
            id="btn-open-audit"
            onClick={onOpenAudit}
            className="flex items-center space-x-1.5 px-2.5 py-1.5 rounded-md text-xs font-medium bg-slate-800 hover:bg-slate-700 border border-slate-700 transition text-slate-200"
            title="Registro de Auditoría (Sección 12)"
          >
            <History className="w-4 h-4 text-amber-400" />
            <span className="hidden sm:inline">Auditoría</span>
            {auditCount > 0 && (
              <span className="px-1.5 py-0.2 rounded-full text-[10px] bg-amber-500/20 text-amber-300 border border-amber-500/30">
                {auditCount}
              </span>
            )}
          </button>

          {/* Go Codebase Explorer */}
          <button
            id="btn-open-codebase"
            onClick={onOpenCodebase}
            className="flex items-center space-x-1.5 px-2.5 py-1.5 rounded-md text-xs font-medium bg-slate-800 hover:bg-slate-700 border border-slate-700 transition text-slate-200"
            title="Explorar código 100% Go Desacoplado"
          >
            <Code className="w-4 h-4 text-cyan-400" />
            <span className="hidden md:inline">Código Go (9 Servicios)</span>
          </button>

          {/* Hybrid View Mode Switch (Desktop vs Mobile Frame) */}
          <button
            id="btn-toggle-hybrid-mode"
            onClick={onToggleMobileView}
            className={`flex items-center space-x-1 px-2.5 py-1.5 rounded-md text-xs font-medium border transition ${
              isMobileView
                ? "bg-purple-950/60 border-purple-500/50 text-purple-300 shadow-sm"
                : "bg-slate-800 hover:bg-slate-700 border-slate-700 text-slate-300"
            }`}
            title="Alternar entre Vista Estación Clínica y Vista App Híbrida Móvil"
          >
            {isMobileView ? (
              <>
                <Monitor className="w-3.5 h-3.5" />
                <span className="hidden sm:inline">Vista Escritorio</span>
              </>
            ) : (
              <>
                <Smartphone className="w-3.5 h-3.5" />
                <span className="hidden sm:inline">App Móvil</span>
              </>
            )}
          </button>
        </div>
      </div>
    </header>
  );
};
