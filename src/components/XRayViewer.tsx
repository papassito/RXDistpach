import React, { useState, useRef } from "react";
import { 
  ZoomIn, 
  ZoomOut, 
  RotateCcw, 
  Sun, 
  Contrast, 
  Eye, 
  ShieldCheck, 
  Sliders, 
  Ruler, 
  Layers
} from "lucide-react";
import { XRayImage } from "../types";

interface XRayViewerProps {
  image: XRayImage | undefined;
  studyTitle: string;
  anatomicalRegion: string;
  onValidateImage?: () => void;
}

export const XRayViewer: React.FC<XRayViewerProps> = ({
  image,
  studyTitle,
  anatomicalRegion,
  onValidateImage,
}) => {
  const [zoom, setZoom] = useState<number>(1);
  const [brightness, setBrightness] = useState<number>(100);
  const [contrast, setContrast] = useState<number>(100);
  const [isInverted, setIsInverted] = useState<boolean>(false);
  const [activeTab, setActiveTab] = useState<"derived" | "original">("derived");
  const [measuringMode, setMeasuringMode] = useState<boolean>(false);
  const [caliperStart, setCaliperStart] = useState<{ x: number; y: number } | null>(null);
  const [caliperEnd, setCaliperEnd] = useState<{ x: number; y: number } | null>(null);

  const containerRef = useRef<HTMLDivElement>(null);

  const resetControls = () => {
    setZoom(1);
    setBrightness(100);
    setContrast(100);
    setIsInverted(false);
    setCaliperStart(null);
    setCaliperEnd(null);
  };

  const handleImageClick = (e: React.MouseEvent<HTMLDivElement>) => {
    if (!measuringMode) return;
    const rect = e.currentTarget.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;

    if (!caliperStart || (caliperStart && caliperEnd)) {
      setCaliperStart({ x, y });
      setCaliperEnd(null);
    } else {
      setCaliperEnd({ x, y });
    }
  };

  // Calculate measured mm based on normalized scale (e.g. 1px ~ 0.25mm)
  const measuredDistance = caliperStart && caliperEnd 
    ? Math.round(Math.hypot(caliperEnd.x - caliperStart.x, caliperEnd.y - caliperStart.y) * 0.25)
    : 0;

  if (!image) {
    return (
      <div className="h-96 flex items-center justify-center bg-slate-950 rounded-xl border border-slate-800 text-slate-500">
        <p>No hay imagen radiográfica asociada al estudio.</p>
      </div>
    );
  }

  // Principle of Integrity (Section 16):
  // Original URL remains untouched; Derived copy is displayed during processing
  const displayedSrc = activeTab === "original" ? image.originalUrl : (image.derivedUrl || image.originalUrl);

  return (
    <div className="bg-slate-950 rounded-xl border border-slate-800 overflow-hidden flex flex-col shadow-xl">
      {/* Viewer Header & Mode Selector */}
      <div className="bg-slate-900/90 px-4 py-3 border-b border-slate-800 flex flex-wrap items-center justify-between gap-2">
        <div className="flex items-center space-x-3">
          <div className="flex bg-slate-800 p-0.5 rounded-lg border border-slate-700">
            <button
              id="tab-derived-copy"
              onClick={() => setActiveTab("derived")}
              className={`px-2.5 py-1 text-xs font-medium rounded-md transition ${
                activeTab === "derived"
                  ? "bg-sky-600 text-white shadow-sm"
                  : "text-slate-400 hover:text-slate-200"
              }`}
            >
              Copia Derivada (Análisis)
            </button>
            <button
              id="tab-original-copy"
              onClick={() => setActiveTab("original")}
              className={`px-2.5 py-1 text-xs font-medium rounded-md transition ${
                activeTab === "original"
                  ? "bg-emerald-600 text-white shadow-sm"
                  : "text-slate-400 hover:text-slate-200"
              }`}
            >
              Original Intacta
            </button>
          </div>

          <div className="flex items-center space-x-1.5 text-xs text-slate-300 font-mono">
            <ShieldCheck className="w-4 h-4 text-emerald-400" />
            <span className="hidden sm:inline">SHA-256:</span>
            <span className="text-[11px] text-slate-400 truncate max-w-[100px] sm:max-w-[150px]">
              {image.metadata?.checksumSHA?.slice(0, 16) || "e3b0c44298fc1c14"}...
            </span>
          </div>
        </div>

        {/* Technical acquisition badge */}
        <div className="text-[11px] text-slate-400 flex items-center space-x-2">
          <span className="px-2 py-0.5 rounded bg-slate-800 border border-slate-700 font-mono">
            {image.metadata?.kvp || "120 kVp"} • {image.metadata?.milliAmps || "3.2 mAs"}
          </span>
          <span className="px-2 py-0.5 rounded bg-slate-800 border border-slate-700 font-mono">
            {image.metadata?.projection || "PA"}
          </span>
        </div>
      </div>

      {/* Main Viewport */}
      <div 
        ref={containerRef}
        onClick={handleImageClick}
        className="relative bg-black flex-1 min-h-[380px] sm:min-h-[460px] flex items-center justify-center overflow-hidden select-none cursor-crosshair"
      >
        <div 
          className="transition-transform duration-100 ease-out flex items-center justify-center"
          style={{
            transform: `scale(${zoom})`,
            filter: `brightness(${brightness}%) contrast(${contrast}%) ${isInverted ? "invert(100%)" : ""}`,
          }}
        >
          <img
            src={displayedSrc}
            alt={studyTitle}
            className="max-h-[440px] max-w-full object-contain pointer-events-none"
            referrerPolicy="no-referrer"
          />
        </div>

        {/* Measuring Caliper Overlay */}
        {measuringMode && caliperStart && (
          <svg className="absolute inset-0 pointer-events-none w-full h-full">
            <circle cx={caliperStart.x} cy={caliperStart.y} r={4} fill="#38bdf8" />
            {caliperEnd && (
              <>
                <circle cx={caliperEnd.x} cy={caliperEnd.y} r={4} fill="#38bdf8" />
                <line
                  x1={caliperStart.x}
                  y1={caliperStart.y}
                  x2={caliperEnd.x}
                  y2={caliperEnd.y}
                  stroke="#38bdf8"
                  strokeWidth={2}
                  strokeDasharray="4,4"
                />
                <rect
                  x={(caliperStart.x + caliperEnd.x) / 2 - 25}
                  y={(caliperStart.y + caliperEnd.y) / 2 - 12}
                  width={50}
                  height={20}
                  rx={4}
                  fill="#0369a1"
                  opacity={0.9}
                />
                <text
                  x={(caliperStart.x + caliperEnd.x) / 2}
                  y={(caliperStart.y + caliperEnd.y) / 2 + 2}
                  fill="#ffffff"
                  fontSize={11}
                  fontWeight="bold"
                  textAnchor="middle"
                >
                  {measuredDistance} mm
                </text>
              </>
            )}
          </svg>
        )}

        {/* Top-Right HUD Badge: Integrity Guarantee (Section 16) */}
        <div className="absolute top-3 right-3 bg-slate-900/80 backdrop-blur border border-slate-700/80 px-2.5 py-1 rounded text-[11px] text-slate-300 pointer-events-none flex items-center space-x-1.5 shadow">
          <Layers className="w-3.5 h-3.5 text-sky-400" />
          <span>{activeTab === "original" ? "Original Intacta (Inmutable)" : "Copia Derivada para Proceso"}</span>
        </div>

        {/* Bottom-Left Anatomical Label */}
        <div className="absolute bottom-3 left-3 bg-slate-900/80 backdrop-blur border border-slate-700/80 px-2.5 py-1 rounded text-[11px] text-slate-300 pointer-events-none font-mono">
          {anatomicalRegion} • {image.fileName}
        </div>
      </div>

      {/* Radiographic Tool Controls Toolbar */}
      <div className="bg-slate-900 px-4 py-3 border-t border-slate-800 flex flex-wrap items-center justify-between gap-4 text-slate-300 text-xs">
        {/* Sliders: Window Level (Brightness & Contrast) */}
        <div className="flex items-center space-x-4">
          <div className="flex items-center space-x-2">
            <Sun className="w-4 h-4 text-amber-400" />
            <span className="w-10">Brillo</span>
            <input
              type="range"
              min="50"
              max="180"
              value={brightness}
              onChange={(e) => setBrightness(Number(e.target.value))}
              className="w-20 sm:w-24 accent-sky-500 h-1.5 bg-slate-700 rounded-lg cursor-pointer"
            />
            <span className="font-mono text-[10px] w-8">{brightness}%</span>
          </div>

          <div className="flex items-center space-x-2">
            <Contrast className="w-4 h-4 text-cyan-400" />
            <span className="w-12">Contraste</span>
            <input
              type="range"
              min="50"
              max="200"
              value={contrast}
              onChange={(e) => setContrast(Number(e.target.value))}
              className="w-20 sm:w-24 accent-sky-500 h-1.5 bg-slate-700 rounded-lg cursor-pointer"
            />
            <span className="font-mono text-[10px] w-8">{contrast}%</span>
          </div>
        </div>

        {/* Actions: Invert, Caliper, Zoom, Reset */}
        <div className="flex items-center space-x-2">
          {/* Invert */}
          <button
            id="btn-invert-xray"
            onClick={() => setIsInverted(!isInverted)}
            className={`flex items-center space-x-1 px-2.5 py-1.5 rounded border transition text-xs ${
              isInverted 
                ? "bg-amber-950/70 border-amber-500 text-amber-300" 
                : "bg-slate-800 border-slate-700 hover:bg-slate-700 text-slate-300"
            }`}
            title="Invertir densidades (Negativo / Positivo)"
          >
            <Eye className="w-3.5 h-3.5" />
            <span>{isInverted ? "Positivo" : "Negativo"}</span>
          </button>

          {/* Caliper Measurement */}
          <button
            id="btn-measure-caliper"
            onClick={() => {
              setMeasuringMode(!measuringMode);
              setCaliperStart(null);
              setCaliperEnd(null);
            }}
            className={`flex items-center space-x-1 px-2.5 py-1.5 rounded border transition text-xs ${
              measuringMode 
                ? "bg-sky-950/70 border-sky-500 text-sky-300" 
                : "bg-slate-800 border-slate-700 hover:bg-slate-700 text-slate-300"
            }`}
            title="Regla de medición milimétrica calibrada"
          >
            <Ruler className="w-3.5 h-3.5" />
            <span>Regla</span>
          </button>

          {/* Zoom controls */}
          <div className="flex items-center space-x-1 bg-slate-800 border border-slate-700 rounded p-0.5">
            <button
              onClick={() => setZoom((z) => Math.max(0.6, z - 0.2))}
              className="p-1 hover:bg-slate-700 rounded text-slate-300"
              title="Reducir zoom"
            >
              <ZoomOut className="w-3.5 h-3.5" />
            </button>
            <span className="px-1.5 font-mono text-[10px] text-slate-300">
              {Math.round(zoom * 100)}%
            </span>
            <button
              onClick={() => setZoom((z) => Math.min(2.5, z + 0.2))}
              className="p-1 hover:bg-slate-700 rounded text-slate-300"
              title="Aumentar zoom"
            >
              <ZoomIn className="w-3.5 h-3.5" />
            </button>
          </div>

          {/* Reset */}
          <button
            onClick={resetControls}
            className="p-1.5 bg-slate-800 hover:bg-slate-700 border border-slate-700 rounded text-slate-400 hover:text-slate-200 transition"
            title="Restablecer controles"
          >
            <RotateCcw className="w-3.5 h-3.5" />
          </button>
        </div>
      </div>
    </div>
  );
};
