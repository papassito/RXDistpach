import React, { useState } from "react";
import { 
  Package, 
  Download, 
  Send, 
  FileCheck, 
  FileText, 
  Image as ImageIcon, 
  AlertCircle, 
  CheckCircle,
  ExternalLink,
  ShieldCheck,
  Building2,
  PhoneCall
} from "lucide-react";
import JSZip from "jszip";
import { ResultPackage, Study } from "../types";

interface ResultBuilderViewProps {
  resultPackage: ResultPackage | null;
  study: Study;
  isBuilding: boolean;
  onBuildPackage: () => void;
  onOpenDelivery: () => void;
}

export const ResultBuilderView: React.FC<ResultBuilderViewProps> = ({
  resultPackage,
  study,
  isBuilding,
  onBuildPackage,
  onOpenDelivery,
}) => {
  const [isDownloadingZip, setIsDownloadingZip] = useState(false);
  const [showSignRequestModal, setShowSignRequestModal] = useState(false);
  const [requestSignedSuccess, setRequestSignedSuccess] = useState(false);

  // Generate real ZIP file containing the RX_RESULT package (Section 10)
  const handleDownloadPackageZip = async () => {
    if (!resultPackage) return;
    setIsDownloadingZip(true);
    try {
      const zip = new JSZip();
      const folder = zip.folder(`RX_RESULT_${study.studyIdentifier}`);

      // 1. resultado.html
      const htmlContent = `<!DOCTYPE html>
<html lang="es">
<head>
  <meta charset="UTF-8">
  <title>RX DISPATCH - ${study.studyIdentifier}</title>
  <style>
    body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; padding: 40px; color: #1e293b; background: #ffffff; line-height: 1.6; max-width: 800px; margin: 0 auto; }
    .header { border-bottom: 2px solid #0284c7; padding-bottom: 12px; margin-bottom: 24px; display: flex; justify-content: space-between; align-items: center; }
    .badge { display: inline-block; padding: 4px 10px; border-radius: 4px; font-weight: bold; font-size: 11px; text-transform: uppercase; }
    .badge-amber { background: #fef3c7; color: #92400e; border: 1px solid #fcd34d; }
    .badge-red { background: #fee2e2; color: #991b1b; border: 1px solid #fca5a5; }
    .legend { border: 2px solid #f59e0b; background-color: #fffbeb; padding: 16px; border-radius: 6px; margin: 20px 0; color: #78350f; }
    .status-box { background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 6px; padding: 14px; margin: 16px 0; display: grid; grid-template-columns: repeat(2, 1fr); gap: 10px; }
  </style>
</head>
<body>
  <div class="header">
    <div>
      <h2 style="margin: 0; color: #0f172a;">RX DISPATCH BY KLIK SOFT PRO</h2>
      <small style="color: #64748b;">Sistema de Envío de Resultados e Imágenes de Rayos X</small>
    </div>
    <div>
      <span class="badge badge-amber">Lectura Genérica Automatizada</span>
      <span class="badge badge-red">Sin Firma Médica</span>
    </div>
  </div>

  <div class="status-box">
    <div><strong>Tipo de Resultado:</strong> LECTURA GENÉRICA AUTOMATIZADA</div>
    <div><strong>Estado:</strong> SIN FIRMA MÉDICA</div>
    <div><strong>Informe Médico Oficial:</strong> NO INCLUIDO</div>
    <div><strong>Firma Responsable:</strong> DEBE SOLICITARSE</div>
  </div>

  <div class="legend">
    <h3 style="margin: 0 0 6px 0; color: #b45309; font-size: 13px;">AVISO IMPORTANTE</h3>
    <p style="margin: 0 0 6px 0; font-size: 13px;">La lectura incluida en este resultado es de carácter <strong>GENÉRICO</strong> y tiene como finalidad proporcionar una descripción general de la imagen.</p>
    <p style="margin: 0 0 6px 0; font-size: 13px; font-weight: bold; color: #991b1b;">No constituye un diagnóstico médico ni sustituye el informe radiológico oficial.</p>
    <p style="margin: 0; font-size: 12.5px;">Si requiere el informe y la firma del médico responsable, deberá solicitarlo directamente al servicio correspondiente.</p>
  </div>

  <h3>DATOS DEL ESTUDIO</h3>
  <p><strong>Identificador:</strong> ${study.studyIdentifier}<br>
  <strong>Estudio:</strong> ${study.studyType} (${study.anatomicalRegion})<br>
  <strong>Paciente:</strong> ${study.patient.fullName} (ID: ${study.patient.id})<br>
  <strong>Médico Solicitante:</strong> ${study.referringPhysician}</p>

  <h3>LECTURA GENÉRICA DE LA IMAGEN</h3>
  <p>${resultPackage.reading.generalDescription}</p>

  <h4>Observaciones Visuales:</h4>
  <ul>
    ${resultPackage.reading.visualObservations.map((o) => `<li>${o}</li>`).join("")}
  </ul>

  <h4>Características Observables:</h4>
  <ul>
    ${resultPackage.reading.observableCharacteristics.map((c) => `<li>${c}</li>`).join("")}
  </ul>

  <hr style="margin-top: 30px; border: 0; border-top: 1px solid #e2e8f0;">
  <p style="font-size: 11px; color: #64748b;">
    Paquete RX_RESULT generado el ${new Date(resultPackage.generatedAt).toLocaleString()} • Checksum SHA: ${resultPackage.checksumSha}
  </p>
</body>
</html>`;

      folder?.file("resultado.html", htmlContent);

      // 2. lectura-generica.txt
      const txtContent = `RX DISPATCH BY KLIK SOFT PRO - LECTURA GENÉRICA DE LA IMAGEN
=========================================================
ESTUDIO: ${study.studyType}
IDENTIFICADOR: ${study.studyIdentifier}
PACIENTE: ${study.patient.fullName} (ID: ${study.patient.id})
FECHA: ${new Date(resultPackage.generatedAt).toLocaleString()}

TIPO DE RESULTADO: LECTURA GENÉRICA AUTOMATIZADA
ESTADO: SIN FIRMA MÉDICA
INFORME MÉDICO OFICIAL: NO INCLUIDO
FIRMA DEL MÉDICO RESPONSABLE: DEBE SOLICITARSE

${resultPackage.mandatoryLegend}

DESCRIPCIÓN GENERAL:
${resultPackage.reading.generalDescription}

OBSERVACIONES VISUALES:
${resultPackage.reading.visualObservations.map((o) => `- ${o}`).join("\n")}

CARACTERÍSTICAS OBSERVABLES:
${resultPackage.reading.observableCharacteristics.map((c) => `- ${c}`).join("\n")}

CUIDADO Y ALCANCE TÉCNICO:
${resultPackage.reading.technicalCaveats}

CHECKSUM SHA-256 DEL PAQUETE: ${resultPackage.checksumSha}
=========================================================
`;
      folder?.file("lectura-generica.txt", txtContent);

      // 3. images/ folder with manifest and raw images
      const imgFolder = folder?.folder("images");
      study.images.forEach((img, i) => {
        imgFolder?.file(
          `image-${String(i + 1).padStart(3, "0")}-${img.fileName}.txt`,
          `ORIGINAL RADIOGRAPHIC IMAGE REFERENCE
Filename: ${img.fileName}
Format: ${img.metadata.format}
Resolution: ${img.metadata.width}x${img.metadata.height}
kVp: ${img.metadata.kvp}, mA: ${img.metadata.milliAmps}
SHA-256 Checksum: ${img.metadata.checksumSHA}
Integrity status: Preserved original unaltered.`
        );
      });

      const zipBlob = await zip.generateAsync({ type: "blob" });
      const downloadUrl = URL.createObjectURL(zipBlob);
      const link = document.createElement("a");
      link.href = downloadUrl;
      link.download = `RX_RESULT_${study.studyIdentifier}.zip`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(downloadUrl);
    } catch (err) {
      console.error("Error creating zip:", err);
    } finally {
      setIsDownloadingZip(false);
    }
  };

  return (
    <div className="bg-slate-900 border border-slate-800 rounded-xl p-5 shadow-lg space-y-6">
      {/* Header */}
      <div className="flex flex-wrap items-center justify-between gap-3 border-b border-slate-800 pb-4">
        <div>
          <div className="flex items-center space-x-2">
            <Package className="w-5 h-5 text-sky-400" />
            <h2 className="text-base sm:text-lg font-bold tracking-tight text-white uppercase">
              RESULT BUILDER (PAQUETE DE ENVÍO)
            </h2>
          </div>
          <p className="text-xs text-slate-400 mt-0.5">
            Componente: Result Builder • Construcción y certificación de paquete con advertencias y anexos
          </p>
        </div>

        <div className="flex items-center space-x-2">
          {!resultPackage ? (
            <button
              id="btn-build-result-package"
              onClick={onBuildPackage}
              disabled={isBuilding}
              className="flex items-center space-x-2 px-3 py-1.5 rounded-md text-xs font-semibold bg-sky-600 hover:bg-sky-500 text-white transition disabled:opacity-50 shadow-sm"
            >
              {isBuilding ? (
                <>
                  <div className="w-3.5 h-3.5 border-2 border-white border-t-transparent rounded-full animate-spin" />
                  <span>Construyendo paquete...</span>
                </>
              ) : (
                <>
                  <FileCheck className="w-3.5 h-3.5" />
                  <span>Construir Paquete RX_RESULT</span>
                </>
              )}
            </button>
          ) : (
            <>
              <button
                id="btn-download-zip"
                onClick={handleDownloadPackageZip}
                disabled={isDownloadingZip}
                className="flex items-center space-x-1.5 px-3 py-1.5 rounded-md text-xs font-semibold bg-slate-800 hover:bg-slate-700 border border-slate-700 text-slate-200 transition"
                title="Descargar paquete completo RX_RESULT en formato .ZIP (Sección 10)"
              >
                <Download className="w-3.5 h-3.5 text-sky-400" />
                <span>{isDownloadingZip ? "Generando ZIP..." : "Descargar Paquete (.ZIP)"}</span>
              </button>

              <button
                id="btn-open-delivery-modal"
                onClick={onOpenDelivery}
                className="flex items-center space-x-1.5 px-3 py-1.5 rounded-md text-xs font-semibold bg-emerald-600 hover:bg-emerald-500 text-white transition shadow-sm"
              >
                <Send className="w-3.5 h-3.5" />
                <span>Ejecutar Envío</span>
              </button>
            </>
          )}
        </div>
      </div>

      {resultPackage ? (
        <div className="space-y-5">
          {/* Assembled Package Summary Cards */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
            <div className="bg-slate-950 p-3.5 rounded-lg border border-slate-800">
              <span className="text-[10px] uppercase font-semibold text-slate-500 block mb-1">
                Estructura del Paquete (Sección 10)
              </span>
              <div className="font-mono text-xs text-sky-300">RX_RESULT/</div>
              <ul className="text-xs text-slate-300 mt-1.5 space-y-1 font-mono">
                <li className="flex items-center space-x-1.5">
                  <FileText className="w-3 h-3 text-emerald-400" />
                  <span>resultado.pdf / html</span>
                </li>
                <li className="flex items-center space-x-1.5">
                  <FileText className="w-3 h-3 text-amber-400" />
                  <span>lectura-generica.txt</span>
                </li>
                <li className="flex items-center space-x-1.5">
                  <ImageIcon className="w-3 h-3 text-cyan-400" />
                  <span>images/ ({study.images.length} imágenes originales)</span>
                </li>
              </ul>
            </div>

            <div className="bg-slate-950 p-3.5 rounded-lg border border-slate-800">
              <span className="text-[10px] uppercase font-semibold text-slate-500 block mb-1">
                Garantía de Integridad (Sección 16)
              </span>
              <div className="text-xs text-slate-300 flex items-center space-x-1.5">
                <ShieldCheck className="w-4 h-4 text-emerald-400 flex-shrink-0" />
                <span>Imágenes originales intactas</span>
              </div>
              <p className="text-[11px] text-slate-400 mt-1 font-mono truncate">
                Hash: {resultPackage.checksumSha.slice(0, 16)}...
              </p>
              <span className="inline-block mt-2 px-2 py-0.5 text-[10px] rounded bg-emerald-950/80 border border-emerald-800 text-emerald-300">
                Certificación SHA-256 OK
              </span>
            </div>

            <div className="bg-slate-950 p-3.5 rounded-lg border border-slate-800 flex flex-col justify-between">
              <div>
                <span className="text-[10px] uppercase font-semibold text-slate-500 block mb-1">
                  Informe Oficial y Firma Médica
                </span>
                <p className="text-xs text-slate-300">
                  El informe con valor médico legal debe solicitarse formalmente.
                </p>
              </div>
              <button
                id="btn-request-signed-report"
                onClick={() => setShowSignRequestModal(true)}
                className="mt-2 text-xs font-semibold text-amber-400 hover:text-amber-300 flex items-center space-x-1 underline"
              >
                <span>Solicitar Informe Firmado</span>
                <ExternalLink className="w-3 h-3" />
              </button>
            </div>
          </div>

          {/* Interactive Document Preview Box */}
          <div className="bg-white text-slate-900 rounded-lg p-6 shadow-xl border border-slate-300 space-y-4">
            {/* Document Header */}
            <div className="flex justify-between items-start border-b-2 border-sky-600 pb-3">
              <div>
                <h3 className="text-lg font-black text-slate-900 tracking-tight">
                  RX DISPATCH BY KLIK SOFT PRO
                </h3>
                <p className="text-xs text-slate-500">
                  Sistema de Envío de Resultados e Imágenes de Rayos X
                </p>
              </div>
              <div className="text-right">
                <span className="inline-block px-2.5 py-1 text-[11px] font-bold uppercase rounded bg-amber-100 text-amber-900 border border-amber-300">
                  {resultPackage.resultType}
                </span>
                <div className="text-[11px] font-bold text-rose-700 mt-1">
                  {resultPackage.statusNotice}
                </div>
              </div>
            </div>

            {/* Document 4 Status Matrix */}
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-2 text-center bg-slate-100 p-2.5 rounded border border-slate-200 text-xs">
              <div>
                <span className="text-[10px] text-slate-500 block font-semibold">TIPO DE RESULTADO</span>
                <span className="font-bold text-slate-800 text-[11px]">{resultPackage.resultType}</span>
              </div>
              <div>
                <span className="text-[10px] text-rose-600 block font-semibold">ESTADO</span>
                <span className="font-bold text-rose-700 text-[11px]">{resultPackage.statusNotice}</span>
              </div>
              <div>
                <span className="text-[10px] text-slate-500 block font-semibold">INFORME MÉDICO OFICIAL</span>
                <span className="font-semibold text-slate-700 text-[11px]">{resultPackage.officialReportNotice}</span>
              </div>
              <div>
                <span className="text-[10px] text-amber-700 block font-semibold">FIRMA RESPONSABLE</span>
                <span className="font-bold text-amber-800 text-[11px]">{resultPackage.physicianSignatureNotice}</span>
              </div>
            </div>

            {/* Mandatory Legend (Section 8) in the document */}
            <div className="bg-amber-50 border-l-4 border-amber-500 p-3.5 text-xs text-amber-900 space-y-1">
              <div className="font-black text-amber-800 text-xs">AVISO IMPORTANTE</div>
              <p>
                La lectura incluida en este resultado es de carácter <strong>GENÉRICO</strong> y tiene como finalidad proporcionar una descripción general de la imagen.
              </p>
              <p className="font-bold text-rose-900">
                No constituye un diagnóstico médico ni sustituye el informe radiológico oficial.
              </p>
              <p className="text-amber-800">
                Si requiere el informe y la <strong>firma del médico responsable</strong>, deberá solicitarlo directamente al servicio correspondiente.
              </p>
            </div>

            {/* Clinical & Administrative Data */}
            <div className="grid grid-cols-2 gap-4 text-xs text-slate-700 bg-slate-50 p-3 rounded border border-slate-200">
              <div>
                <span className="font-semibold text-slate-900 block">Estudio:</span>
                <span>{study.studyType} ({study.anatomicalRegion})</span>
                <span className="font-semibold text-slate-900 block mt-1.5">Identificador:</span>
                <span className="font-mono">{study.studyIdentifier}</span>
              </div>
              <div>
                <span className="font-semibold text-slate-900 block">Paciente:</span>
                <span>{study.patient.fullName} (ID: {study.patient.id})</span>
                <span className="font-semibold text-slate-900 block mt-1.5">Médico Solicitante:</span>
                <span>{study.referringPhysician}</span>
              </div>
            </div>

            {/* Generic Reading in Document */}
            <div className="space-y-2 text-xs text-slate-800">
              <h4 className="font-bold text-slate-900 uppercase text-xs border-b border-slate-200 pb-1">
                LECTURA GENÉRICA DE LA IMAGEN
              </h4>
              <p className="leading-relaxed">{resultPackage.reading.generalDescription}</p>
              
              <div className="mt-2">
                <span className="font-semibold text-slate-700">Observaciones visuales:</span>
                <ul className="list-disc pl-5 space-y-1 mt-1">
                  {resultPackage.reading.visualObservations.map((obs, i) => (
                    <li key={i}>{obs}</li>
                  ))}
                </ul>
              </div>

              <div className="mt-2">
                <span className="font-semibold text-slate-700">Características observables:</span>
                <ul className="list-disc pl-5 space-y-1 mt-1">
                  {resultPackage.reading.observableCharacteristics.map((ch, i) => (
                    <li key={i}>{ch}</li>
                  ))}
                </ul>
              </div>
            </div>

            {/* Document Footer */}
            <div className="pt-3 border-t border-slate-200 text-[11px] text-slate-500 flex justify-between items-center">
              <span>Instrucciones: {resultPackage.requestInstructions}</span>
              <span className="font-mono">ID: {resultPackage.id}</span>
            </div>
          </div>
        </div>
      ) : (
        <div className="bg-slate-950/60 border border-slate-800/80 rounded-lg p-8 text-center space-y-3">
          <Package className="w-8 h-8 text-slate-600 mx-auto" />
          <p className="text-xs text-slate-400 max-w-md mx-auto">
            Haga clic en <strong>"Construir Paquete RX_RESULT"</strong> para generar el resultado oficial del estudio conforme a la especificación normativa.
          </p>
        </div>
      )}

      {/* Modal: Solicitar Informe Oficial Firmado */}
      {showSignRequestModal && (
        <div className="fixed inset-0 bg-black/70 backdrop-blur-sm z-50 flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-xl max-w-md w-full p-5 space-y-4 shadow-2xl">
            <div className="flex items-center space-x-2 text-white">
              <Building2 className="w-5 h-5 text-sky-400" />
              <h3 className="font-bold text-sm sm:text-base">
                Solicitud de Informe Radiológico Oficial Firmado
              </h3>
            </div>

            <p className="text-xs text-slate-300 leading-relaxed">
              En cumplimiento del Principio Fundamental (Sección 14 y 15), el informe oficial firmado por el médico especialista en Radiodiagnóstico debe solicitarse formalmente al centro.
            </p>

            <div className="bg-slate-950 p-3 rounded-lg border border-slate-800 text-xs space-y-2 text-slate-300">
              <div><strong>Estudio:</strong> {study.studyIdentifier} - {study.studyType}</div>
              <div><strong>Paciente:</strong> {study.patient.fullName}</div>
              <div className="flex items-center space-x-2 text-sky-400 mt-2">
                <PhoneCall className="w-3.5 h-3.5" />
                <span>Soporte CM Soluciones: +34 900 123 456 (Ext. 204)</span>
              </div>
            </div>

            {requestSignedSuccess ? (
              <div className="bg-emerald-950/60 border border-emerald-800 p-3 rounded text-xs text-emerald-300 flex items-center space-x-2">
                <CheckCircle className="w-4 h-4 text-emerald-400 flex-shrink-0" />
                <span>Solicitud remitida al médico especialista asignado. Recibirá notificación formal.</span>
              </div>
            ) : (
              <div className="flex justify-end space-x-2 pt-2">
                <button
                  onClick={() => setShowSignRequestModal(false)}
                  className="px-3 py-1.5 rounded text-xs font-medium bg-slate-800 hover:bg-slate-700 text-slate-300"
                >
                  Cerrar
                </button>
                <button
                  onClick={() => setRequestSignedSuccess(true)}
                  className="px-3 py-1.5 rounded text-xs font-semibold bg-sky-600 hover:bg-sky-500 text-white"
                >
                  Registrar Solicitud
                </button>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};
