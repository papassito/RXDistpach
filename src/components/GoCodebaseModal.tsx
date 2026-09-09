import React, { useState, useEffect } from "react";
import { 
  X, 
  Code, 
  Copy, 
  Check, 
  FolderTree, 
  FileCode, 
  Download, 
  ShieldCheck, 
  Layers
} from "lucide-react";
import JSZip from "jszip";
import { GoSourceFile } from "../types";

interface GoCodebaseModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const GoCodebaseModal: React.FC<GoCodebaseModalProps> = ({
  isOpen,
  onClose,
}) => {
  const [files, setFiles] = useState<GoSourceFile[]>([]);
  const [activeFilePath, setActiveFilePath] = useState<string>("");
  const [copied, setCopied] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    if (isOpen && files.length === 0) {
      setIsLoading(true);
      fetch("/api/go-codebase")
        .then((res) => res.json())
        .then((data: GoSourceFile[]) => {
          setFiles(data);
          if (data.length > 0) {
            setActiveFilePath(data[0].path);
          }
        })
        .catch((err) => console.error("Error loading Go codebase:", err))
        .finally(() => setIsLoading(false));
    }
  }, [isOpen, files.length]);

  if (!isOpen) return null;

  const activeFile = files.find((f) => f.path === activeFilePath) || files[0];

  const handleCopy = () => {
    if (activeFile) {
      navigator.clipboard.writeText(activeFile.content);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  const handleDownloadAll = async () => {
    const zip = new JSZip();
    files.forEach((f) => {
      zip.file(f.path, f.content);
    });
    const blob = await zip.generateAsync({ type: "blob" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = "rx-dispatch-100-go-decoupled-mesh.zip";
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  };

  return (
    <div className="fixed inset-0 bg-black/85 backdrop-blur-md z-50 flex items-center justify-center p-3 sm:p-6">
      <div className="bg-slate-900 border border-slate-800 rounded-xl w-full max-w-5xl h-[88vh] flex flex-col shadow-2xl overflow-hidden relative">
        {/* Header */}
        <div className="bg-slate-950 px-6 py-4 border-b border-slate-800 flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <div className="w-8 h-8 rounded-lg bg-emerald-600/20 text-emerald-400 border border-emerald-500/30 flex items-center justify-center">
              <Code className="w-5 h-5" />
            </div>
            <div>
              <div className="flex items-center space-x-2">
                <h3 className="text-base font-bold text-white tracking-tight">
                  RX DISPATCH • ARQUITECTURA 100% GO · CERO MONOLITO
                </h3>
                <span className="text-[11px] px-2 py-0.5 rounded bg-emerald-950/70 border border-emerald-700/60 text-emerald-300 font-mono">
                  9 Binarios Independientes • Go 1.22
                </span>
              </div>
              <p className="text-xs text-slate-400">
                Malla de microservicios autónomos (Gateway, Security, Study, Storage, Image, Reader, Result, Delivery, Audit)
              </p>
            </div>
          </div>

          <div className="flex items-center space-x-2">
            <button
              onClick={handleDownloadAll}
              className="flex items-center space-x-1.5 px-3 py-1.5 rounded-md text-xs font-semibold bg-emerald-600 hover:bg-emerald-500 text-white transition shadow-sm"
              title="Descargar todos los archivos Go en un ZIP"
            >
              <Download className="w-3.5 h-3.5" />
              <span>Descargar Proyecto Go</span>
            </button>
            <button
              onClick={onClose}
              className="text-slate-400 hover:text-white p-1 rounded-md"
            >
              <X className="w-5 h-5" />
            </button>
          </div>
        </div>

        {/* Body: Left Sidebar Tree + Right Editor */}
        <div className="flex-1 flex overflow-hidden">
          {/* File Tree Sidebar */}
          <div className="w-64 sm:w-72 bg-slate-950 border-r border-slate-800 flex flex-col">
            <div className="p-3 border-b border-slate-800/80 text-[11px] font-bold uppercase tracking-wider text-slate-400 flex items-center space-x-1.5">
              <FolderTree className="w-3.5 h-3.5 text-emerald-400" />
              <span>Estructura 100% Go Desacoplada</span>
            </div>

            <div className="flex-1 overflow-y-auto p-2 space-y-1 text-xs font-mono">
              {isLoading ? (
                <div className="p-4 text-slate-500 text-xs text-center">Cargando árbol Go...</div>
              ) : (
                files.map((file) => {
                  const isSelected = file.path === activeFilePath;
                  return (
                    <button
                      key={file.path}
                      onClick={() => setActiveFilePath(file.path)}
                      className={`w-full text-left px-2.5 py-1.5 rounded-md flex items-center space-x-2 transition ${
                        isSelected
                          ? "bg-cyan-950/80 border border-cyan-800/80 text-cyan-200 font-semibold"
                          : "text-slate-400 hover:bg-slate-900 hover:text-slate-200"
                      }`}
                    >
                      <FileCode className={`w-3.5 h-3.5 flex-shrink-0 ${isSelected ? "text-cyan-400" : "text-slate-500"}`} />
                      <span className="truncate text-[11px]">{file.path.replace("rx-dispatch/", "")}</span>
                    </button>
                  );
                })
              )}
            </div>

            {/* Architecture Footer Notice */}
            <div className="p-3 bg-slate-900/60 border-t border-slate-800 text-[10px] text-slate-400 space-y-1">
              <div className="flex items-center space-x-1 text-emerald-400 font-semibold">
                <Layers className="w-3 h-3" />
                <span>Modo Edge (Binario Único)</span>
              </div>
              <p>Compilación que consolida los 9 servicios en un solo ejecutable, conservando sus fronteras lógicas internas.</p>
            </div>
          </div>

          {/* Right: Code Viewer */}
          <div className="flex-1 flex flex-col bg-slate-900 overflow-hidden">
            {activeFile ? (
              <>
                {/* Code Header Bar */}
                <div className="bg-slate-900 px-4 py-2.5 border-b border-slate-800 flex items-center justify-between text-xs">
                  <span className="font-mono text-cyan-300 font-medium">
                    {activeFile.path}
                  </span>
                  <button
                    onClick={handleCopy}
                    className="flex items-center space-x-1.5 px-2.5 py-1 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs border border-slate-700 transition"
                  >
                    {copied ? (
                      <>
                        <Check className="w-3.5 h-3.5 text-emerald-400" />
                        <span className="text-emerald-400">Copiado</span>
                      </>
                    ) : (
                      <>
                        <Copy className="w-3.5 h-3.5" />
                        <span>Copiar Código</span>
                      </>
                    )}
                  </button>
                </div>

                {/* Code Content View */}
                <div className="flex-1 overflow-auto p-4 font-mono text-xs text-slate-200 bg-slate-950/60 leading-relaxed selection:bg-cyan-900">
                  <pre className="whitespace-pre">
                    <code>{activeFile.content}</code>
                  </pre>
                </div>
              </>
            ) : (
              <div className="flex-1 flex items-center justify-center text-slate-500 text-xs">
                Seleccione un archivo para inspeccionar su código fuente.
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};
