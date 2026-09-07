import React, { useState, useEffect } from "react";
import { 
  X, 
  Cpu, 
  CheckCircle2, 
  Terminal, 
  Play, 
  Square, 
  RefreshCw, 
  ShieldCheck, 
  Network, 
  Activity,
  AlertTriangle
} from "lucide-react";

interface GoServiceStatus {
  name: string;
  binary: string;
  port: number;
  compiled: boolean;
  running: boolean;
  statusText: string;
  responsibility: string;
}

interface GoMeshModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const GoMeshModal: React.FC<GoMeshModalProps> = ({ isOpen, onClose }) => {
  const [services, setServices] = useState<GoServiceStatus[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [testOutput, setTestOutput] = useState<string>("");
  const [isTesting, setIsTesting] = useState(false);
  const [actionMessage, setActionMessage] = useState<string>("");

  const loadStatus = () => {
    setIsLoading(true);
    fetch("/api/go-mesh/status")
      .then((res) => res.json())
      .then((data) => {
        if (data.services) {
          setServices(data.services);
        }
      })
      .catch((err) => console.error("Error loading Go mesh status:", err))
      .finally(() => setIsLoading(false));
  };

  useEffect(() => {
    if (isOpen) {
      loadStatus();
    }
  }, [isOpen]);

  const handleRunTests = () => {
    setIsTesting(true);
    setTestOutput("Ejecutando suite de pruebas en Go (go test -v ./tests/...)...");
    fetch("/api/go-mesh/run-tests", { method: "POST" })
      .then((res) => res.json())
      .then((data) => {
        setTestOutput(data.output || "Pruebas finalizadas");
      })
      .catch((err) => {
        setTestOutput(`Error al ejecutar pruebas: ${err}`);
      })
      .finally(() => setIsTesting(false));
  };

  const handleSmokeTest = () => {
    setIsTesting(true);
    setTestOutput("Ejecutando Smoke Test orquestado (scripts/smoke_test.sh)...");
    fetch("/api/go-mesh/smoke-test", { method: "POST" })
      .then((res) => res.json())
      .then((data) => {
        setTestOutput(data.output || "Smoke test finalizado");
      })
      .catch((err) => {
        setTestOutput(`Error al ejecutar smoke test: ${err}`);
      })
      .finally(() => setIsTesting(false));
  };

  const handleStartMesh = () => {
    setActionMessage("Iniciando los 9 binarios Go independientes...");
    fetch("/api/go-mesh/start", { method: "POST" })
      .then((res) => res.json())
      .then(() => {
        setActionMessage("Malla iniciada correctamente.");
        setTimeout(loadStatus, 1500);
      })
      .catch((err) => setActionMessage(`Error al iniciar: ${err}`));
  };

  const handleStopMesh = () => {
    setActionMessage("Deteniendo los procesos Go...");
    fetch("/api/go-mesh/stop", { method: "POST" })
      .then((res) => res.json())
      .then(() => {
        setActionMessage("Malla detenida.");
        setTimeout(loadStatus, 1000);
      })
      .catch((err) => setActionMessage(`Error al detener: ${err}`));
  };

  if (!isOpen) return null;

  const onlineCount = services.filter((s) => s.running).length;

  return (
    <div className="fixed inset-0 bg-black/80 backdrop-blur-sm z-50 flex items-center justify-center p-3 sm:p-6 animate-in fade-in duration-150">
      <div className="bg-slate-900 border border-slate-700/80 rounded-xl w-full max-w-5xl max-h-[90vh] flex flex-col shadow-2xl overflow-hidden ring-1 ring-slate-700/50">
        {/* Header */}
        <div className="bg-slate-950 px-6 py-4 border-b border-slate-800 flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <div className="w-9 h-9 rounded-lg bg-emerald-600/20 text-emerald-400 border border-emerald-500/30 flex items-center justify-center">
              <Network className="w-5 h-5" />
            </div>
            <div>
              <div className="flex items-center space-x-2">
                <h3 className="text-base font-bold text-white tracking-tight">
                  MALLA DE MICROSERVICIOS GO (CERO MONOLITO)
                </h3>
                <span className="text-xs px-2 py-0.5 rounded bg-emerald-950/70 border border-emerald-700/60 text-emerald-300 font-mono">
                  {onlineCount} / {services.length} Online
                </span>
              </div>
              <p className="text-xs text-slate-400">
                9 Binarios Go independientes compilados y comunicados mediante HTTP/JSON REST
              </p>
            </div>
          </div>

          <div className="flex items-center space-x-2">
            <button
              onClick={loadStatus}
              disabled={isLoading}
              className="p-1.5 rounded-md text-slate-400 hover:text-white hover:bg-slate-800 transition"
              title="Actualizar estado"
            >
              <RefreshCw className={`w-4 h-4 ${isLoading ? "animate-spin text-sky-400" : ""}`} />
            </button>
            <button
              onClick={onClose}
              className="text-slate-400 hover:text-white p-1 rounded-md"
            >
              <X className="w-5 h-5" />
            </button>
          </div>
        </div>

        {/* Content Body */}
        <div className="flex-1 overflow-y-auto p-6 space-y-6">
          {/* Controls Bar */}
          <div className="flex flex-wrap items-center justify-between gap-3 p-3 bg-slate-950/60 rounded-lg border border-slate-800">
            <div className="flex items-center space-x-2">
              <button
                onClick={handleStartMesh}
                className="flex items-center space-x-1.5 px-3 py-1.5 rounded-md text-xs font-semibold bg-emerald-600 hover:bg-emerald-500 text-white transition shadow-sm"
              >
                <Play className="w-3.5 h-3.5" />
                <span>Iniciar Malla (start_all.sh)</span>
              </button>
              <button
                onClick={handleStopMesh}
                className="flex items-center space-x-1.5 px-3 py-1.5 rounded-md text-xs font-semibold bg-slate-800 hover:bg-slate-700 text-slate-300 border border-slate-700 transition"
              >
                <Square className="w-3.5 h-3.5" />
                <span>Detener Malla</span>
              </button>
            </div>

            <div className="flex items-center space-x-2">
              <button
                onClick={handleRunTests}
                disabled={isTesting}
                className="flex items-center space-x-1.5 px-3 py-1.5 rounded-md text-xs font-semibold bg-sky-600 hover:bg-sky-500 text-white transition shadow-sm"
              >
                <Terminal className="w-3.5 h-3.5" />
                <span>Ejecutar Tests Go (go test -v)</span>
              </button>
              <button
                onClick={handleSmokeTest}
                disabled={isTesting}
                className="flex items-center space-x-1.5 px-3 py-1.5 rounded-md text-xs font-semibold bg-indigo-600 hover:bg-indigo-500 text-white transition shadow-sm"
              >
                <Activity className="w-3.5 h-3.5" />
                <span>Ejecutar Smoke Test</span>
              </button>
            </div>
          </div>

          {actionMessage && (
            <div className="text-xs px-3 py-2 rounded bg-sky-950/60 border border-sky-800/60 text-sky-300 font-mono">
              {actionMessage}
            </div>
          )}

          {/* Test / Terminal Output Section */}
          {testOutput && (
            <div className="rounded-lg border border-slate-800 bg-slate-950 overflow-hidden">
              <div className="px-3 py-2 bg-slate-900 border-b border-slate-800 text-[11px] font-mono text-slate-300 flex items-center justify-between">
                <div className="flex items-center space-x-2">
                  <Terminal className="w-3.5 h-3.5 text-emerald-400" />
                  <span>Consola de Ejecución Go</span>
                </div>
                <button
                  onClick={() => setTestOutput("")}
                  className="text-slate-500 hover:text-slate-300 text-xs"
                >
                  Limpiar
                </button>
              </div>
              <pre className="p-3 text-xs font-mono text-emerald-300/90 whitespace-pre-wrap max-h-48 overflow-y-auto leading-relaxed bg-black/50">
                {testOutput}
              </pre>
            </div>
          )}

          {/* Grid of 9 Independent Services */}
          <div>
            <div className="text-xs font-bold text-slate-400 uppercase tracking-wider mb-3 flex items-center space-x-2">
              <Cpu className="w-4 h-4 text-emerald-400" />
              <span>Estado de los 9 Componentes Go Autónomos</span>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
              {services.map((svc) => (
                <div
                  key={svc.binary}
                  className="p-3.5 rounded-lg bg-slate-950/70 border border-slate-800 hover:border-slate-700 transition flex flex-col justify-between"
                >
                  <div>
                    <div className="flex items-center justify-between mb-1.5">
                      <span className="font-mono font-bold text-sm text-slate-100">
                        {svc.binary}
                      </span>
                      <span
                        className={`text-[10px] font-mono px-2 py-0.5 rounded border ${
                          svc.running
                            ? "bg-emerald-950/80 text-emerald-300 border-emerald-700/60"
                            : svc.compiled
                            ? "bg-amber-950/80 text-amber-300 border-amber-700/60"
                            : "bg-red-950/80 text-red-400 border-red-800/60"
                        }`}
                      >
                        {svc.statusText}
                      </span>
                    </div>

                    <div className="flex items-center space-x-2 text-xs font-mono text-slate-400 mb-2">
                      <span className="text-sky-400 font-semibold">Port {svc.port}</span>
                      <span>•</span>
                      <span>./cmd/{svc.name}</span>
                    </div>

                    <p className="text-xs text-slate-400 leading-normal">
                      {svc.responsibility}
                    </p>
                  </div>

                  <div className="mt-3 pt-2.5 border-t border-slate-800/80 flex items-center justify-between text-[11px] text-slate-500 font-mono">
                    <span>Health: /healthz</span>
                    {svc.running && (
                      <span className="text-emerald-400 flex items-center space-x-1">
                        <CheckCircle2 className="w-3 h-3 inline" />
                        <span>200 OK</span>
                      </span>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Mandate & Compliance Cards */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pt-2">
            <div className="p-4 rounded-lg bg-emerald-950/20 border border-emerald-800/30">
              <div className="flex items-center space-x-2 text-emerald-400 text-xs font-bold uppercase tracking-wider mb-2">
                <ShieldCheck className="w-4 h-4" />
                <span>Mandato Arquitectónico Resuelto</span>
              </div>
              <ul className="text-xs text-slate-300 space-y-1.5 list-disc list-inside">
                <li><strong className="text-slate-100">100% Go</strong>: Cada servicio es un binario independiente con su propio `main.go`.</li>
                <li><strong className="text-slate-100">Cero Monolítico</strong>: No hay un único ejecutable que contenga todos los servicios.</li>
                <li><strong className="text-slate-100">Contratos Explícitos</strong>: Comunicación mediante DTOs tipados en JSON sobre HTTP REST.</li>
                <li><strong className="text-slate-100">Preservación Inmutable</strong>: Original protegido mediante hash SHA-256 en RX STORAGE.</li>
              </ul>
            </div>

            <div className="p-4 rounded-lg bg-amber-950/20 border border-amber-800/30">
              <div className="flex items-center space-x-2 text-amber-400 text-xs font-bold uppercase tracking-wider mb-2">
                <AlertTriangle className="w-4 h-4" />
                <span>Cumplimiento Legal y Clínico</span>
              </div>
              <p className="text-xs text-slate-300 leading-relaxed">
                El servicio <strong>rx-reader</strong> genera exclusivamente lecturas clasificadas como 
                <code className="text-amber-300 font-mono text-[11px] mx-1">GENERIC_AUTOMATED</code> y 
                <code className="text-amber-300 font-mono text-[11px]">WITHOUT_MEDICAL_SIGNATURE</code>, 
                incorporando inalterablemente la leyenda legal que exige solicitar el informe firmado al médico responsable.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
