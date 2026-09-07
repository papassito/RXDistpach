import React from "react";
import { 
  AlertTriangle, 
  Eye, 
  FileText, 
  ShieldAlert, 
  HelpCircle, 
  Sparkles,
  Layers,
  ArrowRight
} from "lucide-react";
import { GenericReading, Study } from "../types";

interface GenericReadingViewProps {
  reading: GenericReading | null;
  study: Study;
  isGenerating: boolean;
  onGenerateReading: () => void;
  onProceedToBuilder: () => void;
}

export const GenericReadingView: React.FC<GenericReadingViewProps> = ({
  reading,
  study,
  isGenerating,
  onGenerateReading,
  onProceedToBuilder,
}) => {
  return (
    <div className="bg-slate-900 border border-slate-800 rounded-xl p-5 shadow-lg space-y-6">
      {/* Title with Strict Naming Convention (Section 6) */}
      <div className="flex flex-wrap items-center justify-between gap-3 border-b border-slate-800 pb-4">
        <div>
          <div className="flex items-center space-x-2">
            <span className="w-2.5 h-2.5 rounded-full bg-amber-400 animate-pulse" />
            <h2 className="text-base sm:text-lg font-bold tracking-tight text-white uppercase">
              LECTURA GENÉRICA DE LA IMAGEN
            </h2>
          </div>
          <p className="text-xs text-slate-400 mt-0.5">
            Componente: Generic Image Reader • Observación automatizada sin emisión diagnóstica
          </p>
        </div>

        <div className="flex items-center space-x-2">
          <button
            id="btn-trigger-generic-reading"
            onClick={onGenerateReading}
            disabled={isGenerating}
            className="flex items-center space-x-2 px-3 py-1.5 rounded-md text-xs font-semibold bg-amber-600 hover:bg-amber-500 active:bg-amber-700 text-white transition disabled:opacity-50 shadow-sm"
          >
            {isGenerating ? (
              <>
                <div className="w-3.5 h-3.5 border-2 border-white border-t-transparent rounded-full animate-spin" />
                <span>Analizando patrón visual...</span>
              </>
            ) : (
              <>
                <Sparkles className="w-3.5 h-3.5" />
                <span>{reading ? "Regenerar Lectura" : "Iniciar Lectura Genérica"}</span>
              </>
            )}
          </button>

          {reading && (
            <button
              id="btn-proceed-to-builder"
              onClick={onProceedToBuilder}
              className="flex items-center space-x-1.5 px-3 py-1.5 rounded-md text-xs font-semibold bg-sky-600 hover:bg-sky-500 text-white transition shadow-sm"
            >
              <span>Construir Paquete</span>
              <ArrowRight className="w-3.5 h-3.5" />
            </button>
          )}
        </div>
      </div>

      {/* 4 Critical Status Notices (Section 15) */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-2.5">
        <div className="bg-slate-950/70 border border-slate-800 rounded-lg p-2.5 text-center">
          <span className="text-[10px] text-slate-400 font-semibold uppercase tracking-wider block">
            Tipo de Resultado
          </span>
          <span className="text-xs font-bold text-amber-400">
            LECTURA GENÉRICA AUTOMATIZADA
          </span>
        </div>

        <div className="bg-slate-950/70 border border-rose-900/40 rounded-lg p-2.5 text-center">
          <span className="text-[10px] text-rose-400/80 font-semibold uppercase tracking-wider block">
            Estado
          </span>
          <span className="text-xs font-bold text-rose-400">
            SIN FIRMA MÉDICA
          </span>
        </div>

        <div className="bg-slate-950/70 border border-slate-800 rounded-lg p-2.5 text-center">
          <span className="text-[10px] text-slate-400 font-semibold uppercase tracking-wider block">
            Informe Oficial
          </span>
          <span className="text-xs font-semibold text-slate-300">
            NO INCLUIDO
          </span>
        </div>

        <div className="bg-slate-950/70 border border-amber-900/40 rounded-lg p-2.5 text-center">
          <span className="text-[10px] text-amber-400/80 font-semibold uppercase tracking-wider block">
            Firma Responsable
          </span>
          <span className="text-xs font-bold text-amber-300">
            DEBE SOLICITARSE
          </span>
        </div>
      </div>

      {/* Mandatory Legend (Section 8) - Compulsory & Non-optional */}
      <div className="bg-amber-950/40 border-2 border-amber-500/80 rounded-lg p-4 text-amber-200 text-xs shadow-inner">
        <div className="flex items-center space-x-2 font-black text-amber-300 text-xs tracking-wider uppercase mb-1.5">
          <AlertTriangle className="w-4 h-4 text-amber-400 flex-shrink-0" />
          <span>AVISO IMPORTANTE</span>
        </div>
        <p className="leading-relaxed">
          La lectura incluida en este resultado es de carácter <strong>GENÉRICO</strong> y tiene como finalidad proporcionar una descripción general de la imagen.
        </p>
        <p className="font-bold text-amber-100 my-1 leading-relaxed">
          No constituye un diagnóstico médico ni sustituye el informe radiológico oficial.
        </p>
        <p className="leading-relaxed text-amber-300/90">
          Si requiere el informe y la <strong>firma del médico responsable</strong>, deberá solicitarlo directamente al servicio correspondiente.
        </p>
      </div>

      {/* Reading Details Output */}
      {reading ? (
        <div className="space-y-4 bg-slate-950 rounded-lg p-4 border border-slate-800">
          {/* Study Header Information */}
          <div className="text-xs text-slate-300 pb-3 border-b border-slate-800 flex flex-wrap justify-between gap-2">
            <div>
              <span className="text-slate-500 font-semibold uppercase text-[10px] block">Estudio:</span>
              <span className="font-medium text-slate-200">{reading.studyTitle}</span>
            </div>
            <div>
              <span className="text-slate-500 font-semibold uppercase text-[10px] block">Región Anatómica:</span>
              <span className="font-medium text-slate-200">{study.anatomicalRegion}</span>
            </div>
            <div>
              <span className="text-slate-500 font-semibold uppercase text-[10px] block">Paciente:</span>
              <span className="font-medium text-slate-200">{study.patient.fullName}</span>
            </div>
          </div>

          {/* General Description */}
          <div>
            <h4 className="text-xs font-semibold text-slate-300 uppercase tracking-wide flex items-center space-x-1.5 mb-1.5">
              <Eye className="w-3.5 h-3.5 text-sky-400" />
              <span>Descripción General:</span>
            </h4>
            <p className="text-xs text-slate-300 bg-slate-900/60 p-3 rounded border border-slate-800 leading-relaxed font-sans">
              {reading.generalDescription}
            </p>
          </div>

          {/* Visual Observations */}
          <div>
            <h4 className="text-xs font-semibold text-slate-300 uppercase tracking-wide flex items-center space-x-1.5 mb-1.5">
              <FileText className="w-3.5 h-3.5 text-sky-400" />
              <span>Observaciones Visuales:</span>
            </h4>
            <ul className="space-y-1.5 bg-slate-900/60 p-3 rounded border border-slate-800 text-xs text-slate-300">
              {reading.visualObservations.map((obs, idx) => (
                <li key={idx} className="flex items-start space-x-2">
                  <span className="text-sky-400 mt-1 font-bold">•</span>
                  <span>{obs}</span>
                </li>
              ))}
            </ul>
          </div>

          {/* Observable Characteristics */}
          <div>
            <h4 className="text-xs font-semibold text-slate-300 uppercase tracking-wide flex items-center space-x-1.5 mb-1.5">
              <Layers className="w-3.5 h-3.5 text-sky-400" />
              <span>Características Observables:</span>
            </h4>
            <ul className="space-y-1.5 bg-slate-900/60 p-3 rounded border border-slate-800 text-xs text-slate-300">
              {reading.observableCharacteristics.map((char, idx) => (
                <li key={idx} className="flex items-start space-x-2">
                  <span className="text-amber-400 mt-1 font-bold">▪</span>
                  <span>{char}</span>
                </li>
              ))}
            </ul>
          </div>

          {/* Technical Caveat Footer */}
          <div className="pt-2 text-[11px] text-slate-400 border-t border-slate-800/80 italic flex items-center space-x-2">
            <HelpCircle className="w-3.5 h-3.5 text-slate-500 flex-shrink-0" />
            <span>{reading.technicalCaveats}</span>
          </div>
        </div>
      ) : (
        <div className="bg-slate-950/60 border border-slate-800/80 rounded-lg p-8 text-center space-y-3">
          <Eye className="w-8 h-8 text-slate-600 mx-auto" />
          <p className="text-xs text-slate-400 max-w-md mx-auto">
            Haga clic en <strong>"Iniciar Lectura Genérica"</strong> para procesar la copia derivada de la imagen mediante el Generic Image Reader.
          </p>
        </div>
      )}
    </div>
  );
};
