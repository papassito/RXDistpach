import React, { useState, useEffect } from "react";
import { 
  Activity, 
  ChevronDown, 
  Layers, 
  Eye, 
  Package, 
  Send, 
  ShieldCheck, 
  User, 
  Calendar, 
  Tag, 
  FileText, 
  CheckCircle2, 
  RefreshCw,
  Smartphone,
  Info
} from "lucide-react";
import { Header } from "./components/Header";
import { XRayViewer } from "./components/XRayViewer";
import { GenericReadingView } from "./components/GenericReadingView";
import { ResultBuilderView } from "./components/ResultBuilderView";
import { DeliveryModal } from "./components/DeliveryModal";
import { AuditDrawer } from "./components/AuditDrawer";
import { GoCodebaseModal } from "./components/GoCodebaseModal";
import { GoMeshModal } from "./components/GoMeshModal";
import { NewStudyModal } from "./components/NewStudyModal";
import { WorkflowPipeline } from "./components/WorkflowPipeline";
import { 
  Study, 
  GenericReading, 
  ResultPackage, 
  AuditRecord, 
  DeliveryRecord 
} from "./types";

export default function App() {
  const [studies, setStudies] = useState<Study[]>([]);
  const [selectedStudyId, setSelectedStudyId] = useState<string>("");
  const [reading, setReading] = useState<GenericReading | null>(null);
  const [resultPackage, setResultPackage] = useState<ResultPackage | null>(null);
  const [deliveries, setDeliveries] = useState<DeliveryRecord[]>([]);
  const [auditEvents, setAuditEvents] = useState<AuditRecord[]>([]);
  
  // Pipeline active view (0: Visor e Imagen, 1: Lectura Genérica, 2: Result Builder, 3: Entregas)
  const [activeStage, setActiveStage] = useState<number>(0);
  
  // Hybrid App Mode (Desktop vs Phone Viewport)
  const [isMobileView, setIsMobileView] = useState<boolean>(false);

  // Modals & Drawers
  const [isNewStudyModalOpen, setIsNewStudyModalOpen] = useState<boolean>(false);
  const [isAuditDrawerOpen, setIsAuditDrawerOpen] = useState<boolean>(false);
  const [isCodebaseModalOpen, setIsCodebaseModalOpen] = useState<boolean>(false);
  const [isMeshModalOpen, setIsMeshModalOpen] = useState<boolean>(false);
  const [isDeliveryModalOpen, setIsDeliveryModalOpen] = useState<boolean>(false);

  // Loading states
  const [isGeneratingReading, setIsGeneratingReading] = useState<boolean>(false);
  const [isBuildingPackage, setIsBuildingPackage] = useState<boolean>(false);

  // Fetch initial studies and audit events
  const loadStudiesAndAudit = () => {
    fetch("/api/studies")
      .then((res) => res.json())
      .then((data: Study[]) => {
        setStudies(data);
        if (data.length > 0 && !selectedStudyId) {
          setSelectedStudyId(data[0].id);
        }
      })
      .catch((err) => console.error("Error loading studies:", err));

    fetch("/api/audit")
      .then((res) => res.json())
      .then((data: AuditRecord[]) => setAuditEvents(data))
      .catch((err) => console.error("Error loading audit:", err));
  };

  useEffect(() => {
    loadStudiesAndAudit();
  }, []);

  const currentStudy = studies.find((s) => s.id === selectedStudyId) || studies[0];

  // Refresh study specific sub-resources
  useEffect(() => {
    if (currentStudy) {
      // Check if study already had reading/package/deliveries
      fetch(`/api/studies/${currentStudy.id}/deliveries`)
        .then((res) => res.json())
        .then((data) => setDeliveries(data))
        .catch(() => setDeliveries([]));
    }
  }, [currentStudy?.id]);

  // Execute Image Engine validation
  const handleValidateImageEngine = async () => {
    if (!currentStudy) return;
    try {
      const res = await fetch(`/api/studies/${currentStudy.id}/image-engine`, {
        method: "POST",
      });
      const data = await res.json();
      loadStudiesAndAudit();
    } catch (err) {
      console.error("Error validating image:", err);
    }
  };

  // Execute Generic Image Reader (Section 6, 7, 8)
  const handleGenerateReading = async () => {
    if (!currentStudy) return;
    setIsGeneratingReading(true);
    try {
      const res = await fetch(`/api/studies/${currentStudy.id}/generic-reading`, {
        method: "POST",
      });
      const data: GenericReading = await res.json();
      setReading(data);
      loadStudiesAndAudit();
    } catch (err) {
      console.error("Error generating reading:", err);
    } finally {
      setIsGeneratingReading(false);
    }
  };

  // Execute Result Builder (Section 9, 10, 15)
  const handleBuildPackage = async () => {
    if (!currentStudy) return;
    setIsBuildingPackage(true);
    try {
      const res = await fetch(`/api/studies/${currentStudy.id}/build-result`, {
        method: "POST",
      });
      const data: ResultPackage = await res.json();
      setResultPackage(data);
      loadStudiesAndAudit();
    } catch (err) {
      console.error("Error building package:", err);
    } finally {
      setIsBuildingPackage(false);
    }
  };

  // Execute Delivery (Section 11)
  const handleDeliver = async (channel: string, recipient: string): Promise<DeliveryRecord | null> => {
    if (!currentStudy) return null;
    try {
      const res = await fetch(`/api/studies/${currentStudy.id}/deliver`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ channel, recipient }),
      });
      const data = await res.json();
      if (data.success && data.delivery) {
        setDeliveries((prev) => [data.delivery, ...prev]);
        loadStudiesAndAudit();
        return data.delivery;
      }
      return null;
    } catch (err) {
      console.error("Error delivering:", err);
      return null;
    }
  };

  const getStatusBadge = (status?: string) => {
    switch (status) {
      case "RECEIVED":
        return "bg-sky-500/20 text-sky-300 border-sky-500/30";
      case "PROCESSING":
        return "bg-indigo-500/20 text-indigo-300 border-indigo-500/30";
      case "GENERIC_READ_COMPLETED":
        return "bg-amber-500/20 text-amber-300 border-amber-500/30";
      case "PACKAGE_BUILT":
        return "bg-cyan-500/20 text-cyan-300 border-cyan-500/30";
      case "DELIVERED":
        return "bg-emerald-500/20 text-emerald-300 border-emerald-500/30";
      default:
        return "bg-slate-800 text-slate-400 border-slate-700";
    }
  };

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans selection:bg-sky-500/30">
      {/* Top Application Header */}
      <Header
        isMobileView={isMobileView}
        onToggleMobileView={() => setIsMobileView(!isMobileView)}
        onOpenNewStudy={() => setIsNewStudyModalOpen(true)}
        onOpenAudit={() => setIsAuditDrawerOpen(true)}
        onOpenCodebase={() => setIsCodebaseModalOpen(true)}
        onOpenMesh={() => setIsMeshModalOpen(true)}
        auditCount={auditEvents.length}
      />

      {/* Main Content Area: Responsive Hybrid Frame or Full Desktop */}
      <main className={`flex-1 ${isMobileView ? "py-6 px-4 flex items-center justify-center bg-slate-900/50" : "max-w-7xl w-full mx-auto p-4 sm:p-6 lg:p-8"}`}>
        {isMobileView ? (
          /* Mobile Hybrid Frame Viewport (Simulating clinician / patient smartphone) */
          <div className="w-[390px] h-[820px] bg-slate-950 border-[10px] border-slate-800 rounded-[44px] shadow-2xl overflow-hidden flex flex-col relative ring-1 ring-slate-700/50">
            {/* Phone Notch / Dynamic Island */}
            <div className="bg-slate-950 h-7 w-full flex items-center justify-center pt-1 flex-shrink-0">
              <div className="w-24 h-4 bg-black rounded-full border border-slate-800" />
            </div>

            {/* Mobile App Bar */}
            <div className="bg-slate-900 px-4 py-2.5 border-b border-slate-800 flex items-center justify-between flex-shrink-0">
              <div className="flex items-center space-x-2">
                <Activity className="w-4 h-4 text-sky-400" />
                <span className="font-bold text-xs tracking-tight">RX DISPATCH APP</span>
              </div>
              <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-amber-950/80 text-amber-400 border border-amber-800">
                Go-Zero Client
              </span>
            </div>

            {/* Mobile Scrollable Viewport */}
            <div className="flex-1 overflow-y-auto p-3.5 space-y-4">
              {/* Study Selector */}
              <div className="bg-slate-900 p-3 rounded-xl border border-slate-800 text-xs">
                <div className="flex justify-between items-center mb-2">
                  <span className="text-[10px] font-semibold text-slate-400 uppercase">Estudio Actual</span>
                  <span className={`px-2 py-0.5 rounded text-[10px] font-mono border ${getStatusBadge(currentStudy?.status)}`}>
                    {currentStudy?.status || "RECEIVED"}
                  </span>
                </div>
                <select
                  value={selectedStudyId}
                  onChange={(e) => setSelectedStudyId(e.target.value)}
                  className="w-full bg-slate-950 border border-slate-700 rounded-md p-2 text-xs text-white focus:outline-none"
                >
                  {studies.map((s) => (
                    <option key={s.id} value={s.id}>
                      {s.studyIdentifier} - {s.patient.fullName} ({s.anatomicalRegion})
                    </option>
                  ))}
                </select>
              </div>

              {/* Radiographic Viewer */}
              {currentStudy && (
                <XRayViewer
                  image={currentStudy.images[0]}
                  studyTitle={currentStudy.studyType}
                  anatomicalRegion={currentStudy.anatomicalRegion}
                  onValidateImage={handleValidateImageEngine}
                />
              )}

              {/* Generic Image Reader */}
              {currentStudy && (
                <GenericReadingView
                  reading={reading}
                  study={currentStudy}
                  isGenerating={isGeneratingReading}
                  onGenerateReading={handleGenerateReading}
                  onProceedToBuilder={() => setActiveStage(2)}
                />
              )}

              {/* Result Builder */}
              {currentStudy && (
                <ResultBuilderView
                  resultPackage={resultPackage}
                  study={currentStudy}
                  isBuilding={isBuildingPackage}
                  onBuildPackage={handleBuildPackage}
                  onOpenDelivery={() => setIsDeliveryModalOpen(true)}
                />
              )}
            </div>

            {/* Mobile Bottom Navigation Bar */}
            <div className="bg-slate-900 border-t border-slate-800 px-4 py-2 flex justify-around text-slate-400 text-[10px] flex-shrink-0">
              <button 
                onClick={() => setActiveStage(0)}
                className="flex flex-col items-center space-y-0.5 text-sky-400"
              >
                <Layers className="w-4 h-4" />
                <span>Imagen</span>
              </button>
              <button 
                onClick={() => {
                  setActiveStage(1);
                  if (!reading) handleGenerateReading();
                }}
                className="flex flex-col items-center space-y-0.5 hover:text-slate-200"
              >
                <Eye className="w-4 h-4" />
                <span>Lectura</span>
              </button>
              <button 
                onClick={() => {
                  setActiveStage(2);
                  if (!resultPackage) handleBuildPackage();
                }}
                className="flex flex-col items-center space-y-0.5 hover:text-slate-200"
              >
                <Package className="w-4 h-4" />
                <span>Paquete</span>
              </button>
              <button 
                onClick={() => setIsDeliveryModalOpen(true)}
                className="flex flex-col items-center space-y-0.5 hover:text-slate-200"
              >
                <Send className="w-4 h-4" />
                <span>Envío</span>
              </button>
            </div>
          </div>
        ) : (
          /* Full-Screen Clinical Desktop Station */
          <div className="space-y-6">
            {/* Top Row: Study Card & Workflow Stepper */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
              {/* Active Study Information Card */}
              <div className="bg-slate-900 border border-slate-800 rounded-xl p-4 shadow-md flex flex-col justify-between">
                <div>
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-[10px] font-bold uppercase tracking-wider text-slate-400">
                      Estudio Seleccionado
                    </span>
                    <span className={`px-2 py-0.5 rounded text-[11px] font-mono border ${getStatusBadge(currentStudy?.status)}`}>
                      {currentStudy?.status || "RECEIVED"}
                    </span>
                  </div>

                  {/* Study Selector Dropdown */}
                  <div className="relative mb-3">
                    <select
                      value={selectedStudyId}
                      onChange={(e) => {
                        setSelectedStudyId(e.target.value);
                        setReading(null);
                        setResultPackage(null);
                      }}
                      className="w-full bg-slate-950 border border-slate-700 rounded-lg py-2 pl-3 pr-8 text-xs font-semibold text-white focus:outline-none focus:border-sky-500 appearance-none cursor-pointer"
                    >
                      {studies.map((s) => (
                        <option key={s.id} value={s.id}>
                          {s.studyIdentifier} • {s.patient.fullName} ({s.anatomicalRegion})
                        </option>
                      ))}
                    </select>
                    <ChevronDown className="w-4 h-4 text-slate-400 absolute right-2.5 top-2.5 pointer-events-none" />
                  </div>

                  {/* Patient Details */}
                  {currentStudy && (
                    <div className="space-y-1.5 text-xs text-slate-300">
                      <div className="flex items-center space-x-2">
                        <User className="w-3.5 h-3.5 text-sky-400" />
                        <span className="font-semibold">{currentStudy.patient.fullName}</span>
                        <span className="text-[11px] text-slate-500 font-mono">({currentStudy.patient.id})</span>
                      </div>
                      <div className="flex items-center space-x-2">
                        <Tag className="w-3.5 h-3.5 text-amber-400" />
                        <span>{currentStudy.studyType}</span>
                      </div>
                      <div className="text-[11px] text-slate-400 italic mt-1 line-clamp-2">
                        "{currentStudy.clinicalIndication}"
                      </div>
                    </div>
                  )}
                </div>

                <div className="pt-3 border-t border-slate-800 flex justify-between items-center text-[11px] text-slate-500">
                  <span>Médico: {currentStudy?.referringPhysician.split("(")[0]}</span>
                  <button
                    onClick={loadStudiesAndAudit}
                    className="hover:text-slate-300 transition p-1"
                    title="Actualizar lista"
                  >
                    <RefreshCw className="w-3.5 h-3.5" />
                  </button>
                </div>
              </div>

              {/* 5-Step Pipeline Stepper */}
              <div className="lg:col-span-2 flex flex-col justify-center">
                <WorkflowPipeline
                  currentStatus={currentStudy?.status || "RECEIVED"}
                  activeStage={activeStage}
                  onSelectStage={(idx) => setActiveStage(idx)}
                />

                {/* Architecture Banner */}
                <div className="mt-3 bg-slate-900/60 border border-slate-800/80 rounded-lg px-4 py-2.5 flex items-center justify-between text-xs text-slate-400">
                  <div className="flex items-center space-x-2">
                    <ShieldCheck className="w-4 h-4 text-emerald-400" />
                    <span>
                      <strong>Separación de Responsabilidades:</strong> Imagen &rarr; Lectura Genérica &rarr; Informe Médico Firmado.
                    </span>
                  </div>
                  <span className="text-[11px] font-mono text-slate-500 hidden sm:inline">
                    Principio Fundamental (Sec. 14)
                  </span>
                </div>
              </div>
            </div>

            {/* Main Stage View Panels */}
            <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
              {/* Left Column (Viewer): 7 Columns on Large Screens */}
              <div className="lg:col-span-7 space-y-4">
                {currentStudy && (
                  <XRayViewer
                    image={currentStudy.images[0]}
                    studyTitle={currentStudy.studyType}
                    anatomicalRegion={currentStudy.anatomicalRegion}
                    onValidateImage={handleValidateImageEngine}
                  />
                )}
              </div>

              {/* Right Column: Workflow Steps & Output (5 Columns) */}
              <div className="lg:col-span-5 space-y-5">
                {/* Generic Image Reader View */}
                {currentStudy && (
                  <GenericReadingView
                    reading={reading}
                    study={currentStudy}
                    isGenerating={isGeneratingReading}
                    onGenerateReading={handleGenerateReading}
                    onProceedToBuilder={() => {
                      setActiveStage(3);
                      handleBuildPackage();
                    }}
                  />
                )}

                {/* Result Builder & Package View */}
                {currentStudy && (
                  <ResultBuilderView
                    resultPackage={resultPackage}
                    study={currentStudy}
                    isBuilding={isBuildingPackage}
                    onBuildPackage={handleBuildPackage}
                    onOpenDelivery={() => setIsDeliveryModalOpen(true)}
                  />
                )}
              </div>
            </div>
          </div>
        )}
      </main>

      {/* Footer Notice */}
      <footer className="bg-slate-900 border-t border-slate-800 py-3 text-center text-xs text-slate-500">
        <p>
          RX DISPATCH BY KLIK • Sistema de envío de resultados e imágenes de Rayos X • Arquitectura 100% Go Desacoplada (9 Microservicios)
        </p>
      </footer>

      {/* Delivery Engine Modal */}
      {currentStudy && (
        <DeliveryModal
          isOpen={isDeliveryModalOpen}
          onClose={() => setIsDeliveryModalOpen(false)}
          study={currentStudy}
          resultPackage={resultPackage}
          onDeliver={handleDeliver}
        />
      )}

      {/* Audit Drawer (Section 12) */}
      <AuditDrawer
        isOpen={isAuditDrawerOpen}
        onClose={() => setIsAuditDrawerOpen(false)}
        auditEvents={auditEvents}
      />

      {/* Go-Zero Codebase Modal */}
      <GoCodebaseModal
        isOpen={isCodebaseModalOpen}
        onClose={() => setIsCodebaseModalOpen(false)}
      />

      {/* Go Mesh Monitoring Modal */}
      <GoMeshModal
        isOpen={isMeshModalOpen}
        onClose={() => setIsMeshModalOpen(false)}
      />

      {/* New Study Intake Modal */}
      <NewStudyModal
        isOpen={isNewStudyModalOpen}
        onClose={() => setIsNewStudyModalOpen(false)}
        onStudyCreated={(newStudy) => {
          setStudies((prev) => [newStudy, ...prev]);
          setSelectedStudyId(newStudy.id);
          setReading(null);
          setResultPackage(null);
          setActiveStage(0);
        }}
      />
    </div>
  );
}
