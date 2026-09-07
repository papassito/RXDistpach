import React, { useState } from "react";
import { 
  X, 
  UploadCloud, 
  FilePlus, 
  Activity, 
  User, 
  FileText, 
  CheckCircle 
} from "lucide-react";
import { Study } from "../types";

interface NewStudyModalProps {
  isOpen: boolean;
  onClose: () => void;
  onStudyCreated: (newStudy: Study) => void;
}

export const NewStudyModal: React.FC<NewStudyModalProps> = ({
  isOpen,
  onClose,
  onStudyCreated,
}) => {
  const [patientName, setPatientName] = useState("");
  const [patientId, setPatientId] = useState("");
  const [anatomicalRegion, setAnatomicalRegion] = useState("Tórax");
  const [studyType, setStudyType] = useState("Radiografía de Tórax PA");
  const [physician, setPhysician] = useState("Dra. Beatriz Morales");
  const [indication, setIndication] = useState("Evaluación radiográfica preventiva de control.");
  const [imageDataUrl, setImageDataUrl] = useState<string>("");
  const [fileName, setFileName] = useState<string>("");
  const [presetSelected, setPresetSelected] = useState<"chest" | "wrist" | "custom">("chest");

  if (!isOpen) return null;

  const handleFileUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setFileName(file.name);
      setPresetSelected("custom");
      const reader = new FileReader();
      reader.onload = () => {
        setImageDataUrl(reader.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  const handlePresetSelect = (preset: "chest" | "wrist") => {
    setPresetSelected(preset);
    if (preset === "chest") {
      setAnatomicalRegion("Tórax");
      setStudyType("Radiografía de Tórax PA");
      setImageDataUrl("/assets/samples/chest_xray_pa.svg");
      setFileName("torax_pa_sample.svg");
    } else {
      setAnatomicalRegion("Muñeca / Extremidad");
      setStudyType("Radiografía de Muñeca AP y Lateral");
      setImageDataUrl("/assets/samples/wrist_xray.svg");
      setFileName("muneca_ap_sample.svg");
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const selectedImage = imageDataUrl || (presetSelected === "chest" ? "/assets/samples/chest_xray_pa.svg" : "/assets/samples/wrist_xray.svg");
    const selectedFileName = fileName || (presetSelected === "chest" ? "torax_pa_sample.svg" : "muneca_ap_sample.svg");

    const studyPayload: Partial<Study> = {
      studyIdentifier: `RX-2026-${Math.floor(1000 + Math.random() * 9000)}`,
      studyType,
      anatomicalRegion,
      referringPhysician: physician,
      clinicalIndication: indication,
      patient: {
        id: patientId || `PAC-${Math.floor(1000 + Math.random() * 9000)}`,
        fullName: patientName || "Paciente Sin Nombre",
        birthDate: "1990-01-01",
        gender: "No especificado",
        email: "paciente@hospital.com",
        phone: "+34 600 123 456",
      },
      images: [
        {
          id: `IMG-${Date.now()}-1`,
          studyId: "",
          fileName: selectedFileName,
          originalUrl: selectedImage,
          derivedUrl: selectedImage,
          createdAt: new Date().toISOString(),
          isIntegrityOk: true,
          metadata: {
            format: "DICOM-Derived / PNG",
            width: 2048,
            height: 2048,
            bitDepth: 16,
            kvp: anatomicalRegion.includes("Tórax") ? "120 kVp" : "55 kVp",
            milliAmps: anatomicalRegion.includes("Tórax") ? "3.2 mAs" : "4.0 mAs",
            projection: anatomicalRegion.includes("Tórax") ? "PA" : "AP",
            checksumSHA: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
          },
        },
      ],
    };

    fetch("/api/studies", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(studyPayload),
    })
      .then((res) => res.json())
      .then((data: Study) => {
        onStudyCreated(data);
        onClose();
      })
      .catch((err) => console.error("Error creating study:", err));
  };

  return (
    <div className="fixed inset-0 bg-black/80 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div className="bg-slate-900 border border-slate-800 rounded-xl max-w-lg w-full p-6 space-y-4 shadow-2xl relative">
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-slate-400 hover:text-white p-1 rounded-md"
        >
          <X className="w-5 h-5" />
        </button>

        <div className="flex items-center space-x-2">
          <FilePlus className="w-5 h-5 text-sky-400" />
          <h3 className="text-base font-bold text-white uppercase tracking-tight">
            RECEPCIÓN DE ESTUDIO RADIOGRÁFICO
          </h3>
        </div>
        <p className="text-xs text-slate-400">
          Componente: Reception & Ingestion • Entrada de estudio radiográfico al pipeline RX Dispatch
        </p>

        <form onSubmit={handleSubmit} className="space-y-4 text-xs text-slate-300">
          {/* Patient Info */}
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="block font-semibold mb-1">Nombre del Paciente:</label>
              <input
                type="text"
                required
                value={patientName}
                onChange={(e) => setPatientName(e.target.value)}
                placeholder="Ej. Roberto Díaz Cano"
                className="w-full bg-slate-950 border border-slate-800 rounded px-2.5 py-1.5 text-white focus:outline-none focus:border-sky-500"
              />
            </div>
            <div>
              <label className="block font-semibold mb-1">ID Paciente / DNI:</label>
              <input
                type="text"
                value={patientId}
                onChange={(e) => setPatientId(e.target.value)}
                placeholder="Ej. PAC-9921"
                className="w-full bg-slate-950 border border-slate-800 rounded px-2.5 py-1.5 text-white focus:outline-none focus:border-sky-500 font-mono"
              />
            </div>
          </div>

          {/* Anatomical Region & Study Type */}
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="block font-semibold mb-1">Región Anatómica:</label>
              <select
                value={anatomicalRegion}
                onChange={(e) => {
                  setAnatomicalRegion(e.target.value);
                  setStudyType(`Radiografía de ${e.target.value} PA/AP`);
                }}
                className="w-full bg-slate-950 border border-slate-800 rounded px-2.5 py-1.5 text-white focus:outline-none focus:border-sky-500"
              >
                <option value="Tórax">Tórax</option>
                <option value="Muñeca">Muñeca</option>
                <option value="Columna Lumbar">Columna Lumbar</option>
                <option value="Abdomen Simple">Abdomen Simple</option>
                <option value="Rodilla">Rodilla</option>
              </select>
            </div>
            <div>
              <label className="block font-semibold mb-1">Tipo de Estudio:</label>
              <input
                type="text"
                value={studyType}
                onChange={(e) => setStudyType(e.target.value)}
                className="w-full bg-slate-950 border border-slate-800 rounded px-2.5 py-1.5 text-white focus:outline-none focus:border-sky-500"
              />
            </div>
          </div>

          {/* Preset Images or Custom Upload */}
          <div>
            <label className="block font-semibold mb-1.5">Imagen Radiográfica:</label>
            <div className="grid grid-cols-2 gap-2 mb-2">
              <button
                type="button"
                onClick={() => handlePresetSelect("chest")}
                className={`p-2 rounded border text-left transition ${
                  presetSelected === "chest"
                    ? "bg-sky-950/70 border-sky-500 text-sky-200"
                    : "bg-slate-950 border-slate-800 text-slate-400"
                }`}
              >
                <div className="font-bold">Muestra: Tórax PA</div>
                <div className="text-[10px] text-slate-500">120 kVp, 3.2 mAs</div>
              </button>

              <button
                type="button"
                onClick={() => handlePresetSelect("wrist")}
                className={`p-2 rounded border text-left transition ${
                  presetSelected === "wrist"
                    ? "bg-sky-950/70 border-sky-500 text-sky-200"
                    : "bg-slate-950 border-slate-800 text-slate-400"
                }`}
              >
                <div className="font-bold">Muestra: Muñeca AP</div>
                <div className="text-[10px] text-slate-500">55 kVp, 4.0 mAs</div>
              </button>
            </div>

            {/* Custom file upload */}
            <label className="border border-dashed border-slate-700 hover:border-sky-500 rounded-lg p-3 flex flex-col items-center justify-center cursor-pointer bg-slate-950/50 hover:bg-slate-950 transition">
              <UploadCloud className="w-5 h-5 text-sky-400 mb-1" />
              <span className="text-slate-300 font-medium">
                {fileName ? fileName : "O sube tu propia placa radiográfica (PNG, JPEG, DICOM)"}
              </span>
              <span className="text-[10px] text-slate-500 mt-0.5">
                La imagen original se conserva intacta bajo el principio de integridad.
              </span>
              <input
                type="file"
                accept="image/*"
                onChange={handleFileUpload}
                className="hidden"
              />
            </label>
          </div>

          <div className="flex justify-end space-x-2 pt-2 border-t border-slate-800">
            <button
              type="button"
              onClick={onClose}
              className="px-3 py-1.5 rounded bg-slate-800 hover:bg-slate-700 text-slate-300 font-medium"
            >
              Cancelar
            </button>
            <button
              type="submit"
              className="px-4 py-1.5 rounded bg-sky-600 hover:bg-sky-500 text-white font-semibold shadow-sm"
            >
              Ingresar Estudio
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
