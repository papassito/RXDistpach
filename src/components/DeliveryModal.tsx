import React, { useState } from "react";
import { 
  Send, 
  Mail, 
  MessageSquare, 
  Globe, 
  CheckCircle2, 
  X, 
  ShieldCheck, 
  Clock, 
  Copy,
  ExternalLink
} from "lucide-react";
import { DeliveryRecord, ResultPackage, Study } from "../types";

interface DeliveryModalProps {
  isOpen: boolean;
  onClose: () => void;
  study: Study;
  resultPackage: ResultPackage | null;
  onDeliver: (channel: string, recipient: string) => Promise<DeliveryRecord | null>;
}

export const DeliveryModal: React.FC<DeliveryModalProps> = ({
  isOpen,
  onClose,
  study,
  resultPackage,
  onDeliver,
}) => {
  const [channel, setChannel] = useState<"EMAIL" | "WHATSAPP" | "SECURE_PORTAL">("EMAIL");
  const [recipient, setRecipient] = useState<string>(study.patient.email || "");
  const [isDelivering, setIsDelivering] = useState(false);
  const [deliveryResult, setDeliveryResult] = useState<DeliveryRecord | null>(null);
  const [copiedToken, setCopiedToken] = useState(false);

  if (!isOpen) return null;

  const handleChannelSelect = (newChannel: "EMAIL" | "WHATSAPP" | "SECURE_PORTAL") => {
    setChannel(newChannel);
    if (newChannel === "EMAIL") {
      setRecipient(study.patient.email || "paciente@hospital.com");
    } else if (newChannel === "WHATSAPP") {
      setRecipient(study.patient.phone || "+34 612 345 678");
    } else {
      setRecipient(`https://rxdispatch.solsuol.net/portal/${study.studyIdentifier}`);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsDelivering(true);
    try {
      const rec = await onDeliver(channel, recipient);
      setDeliveryResult(rec);
    } finally {
      setIsDelivering(false);
    }
  };

  const copyToken = (token: string) => {
    navigator.clipboard.writeText(token);
    setCopiedToken(true);
    setTimeout(() => setCopiedToken(false), 2000);
  };

  return (
    <div className="fixed inset-0 bg-black/80 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div className="bg-slate-900 border border-slate-800 rounded-xl max-w-lg w-full p-6 space-y-5 shadow-2xl relative">
        {/* Close Button */}
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-slate-400 hover:text-white p-1 rounded-md"
        >
          <X className="w-5 h-5" />
        </button>

        {/* Modal Title */}
        <div>
          <div className="flex items-center space-x-2">
            <Send className="w-5 h-5 text-emerald-400" />
            <h3 className="text-base sm:text-lg font-bold text-white uppercase tracking-tight">
              DELIVERY ENGINE - DESPACHO AUTORIZADO
            </h3>
          </div>
          <p className="text-xs text-slate-400 mt-1">
            Sección 11 • Envío seguro del paquete RX_RESULT (resultado + imágenes originales)
          </p>
        </div>

        {deliveryResult ? (
          /* Delivery Success Receipt */
          <div className="space-y-4 bg-slate-950 p-5 rounded-lg border border-slate-800">
            <div className="flex items-center space-x-2 text-emerald-400 font-bold text-sm">
              <CheckCircle2 className="w-5 h-5" />
              <span>Despacho Ejecutado Exitosamente</span>
            </div>

            <div className="space-y-2 text-xs text-slate-300">
              <div className="flex justify-between py-1 border-b border-slate-800">
                <span className="text-slate-500">Estado de Entrega:</span>
                <span className="font-bold text-emerald-400 font-mono">{deliveryResult.status}</span>
              </div>
              <div className="flex justify-between py-1 border-b border-slate-800">
                <span className="text-slate-500">Canal:</span>
                <span className="font-semibold">{deliveryResult.channel}</span>
              </div>
              <div className="flex justify-between py-1 border-b border-slate-800">
                <span className="text-slate-500">Destinatario:</span>
                <span className="font-mono text-slate-200">{deliveryResult.recipient}</span>
              </div>
              <div className="flex justify-between py-1 border-b border-slate-800">
                <span className="text-slate-500">Fecha/Hora:</span>
                <span>{deliveryResult.sentAt ? new Date(deliveryResult.sentAt).toLocaleString() : "Justo ahora"}</span>
              </div>
              <div className="flex justify-between items-center py-2 bg-slate-900/80 px-2.5 rounded border border-slate-700/60 mt-2">
                <div>
                  <span className="text-[10px] text-slate-400 uppercase font-semibold block">Token de Seguimiento:</span>
                  <span className="font-mono text-sky-400 font-bold text-sm">{deliveryResult.trackingToken}</span>
                </div>
                <button
                  onClick={() => copyToken(deliveryResult.trackingToken)}
                  className="px-2 py-1 rounded bg-slate-800 hover:bg-slate-700 text-[11px] text-slate-200 flex items-center space-x-1"
                >
                  <Copy className="w-3 h-3" />
                  <span>{copiedToken ? "Copiado" : "Copiar"}</span>
                </button>
              </div>
            </div>

            <p className="text-[11px] text-slate-400 italic">
              El destinatario recibe el paquete sellado con la leyenda obligatoria y acceso a las imágenes originales no alteradas.
            </p>

            <button
              onClick={onClose}
              className="w-full py-2 bg-slate-800 hover:bg-slate-700 rounded-md text-xs font-semibold text-white transition"
            >
              Cerrar Despachador
            </button>
          </div>
        ) : (
          /* Form for Dispatch */
          <form onSubmit={handleSubmit} className="space-y-4">
            {/* Channel Selection */}
            <div>
              <label className="text-xs font-semibold text-slate-300 block mb-2">
                Seleccione Canal de Entrega Autorizado:
              </label>
              <div className="grid grid-cols-3 gap-2">
                <button
                  type="button"
                  onClick={() => handleChannelSelect("EMAIL")}
                  className={`p-3 rounded-lg border text-xs flex flex-col items-center space-y-1.5 transition ${
                    channel === "EMAIL"
                      ? "bg-sky-950/70 border-sky-500 text-sky-300 font-bold"
                      : "bg-slate-950 border-slate-800 text-slate-400 hover:text-slate-200"
                  }`}
                >
                  <Mail className="w-4 h-4" />
                  <span>Correo Seguro</span>
                </button>

                <button
                  type="button"
                  onClick={() => handleChannelSelect("WHATSAPP")}
                  className={`p-3 rounded-lg border text-xs flex flex-col items-center space-y-1.5 transition ${
                    channel === "WHATSAPP"
                      ? "bg-emerald-950/70 border-emerald-500 text-emerald-300 font-bold"
                      : "bg-slate-950 border-slate-800 text-slate-400 hover:text-slate-200"
                  }`}
                >
                  <MessageSquare className="w-4 h-4" />
                  <span>WhatsApp / SMS</span>
                </button>

                <button
                  type="button"
                  onClick={() => handleChannelSelect("SECURE_PORTAL")}
                  className={`p-3 rounded-lg border text-xs flex flex-col items-center space-y-1.5 transition ${
                    channel === "SECURE_PORTAL"
                      ? "bg-purple-950/70 border-purple-500 text-purple-300 font-bold"
                      : "bg-slate-950 border-slate-800 text-slate-400 hover:text-slate-200"
                  }`}
                >
                  <Globe className="w-4 h-4" />
                  <span>Portal Paciente</span>
                </button>
              </div>
            </div>

            {/* Recipient Input */}
            <div>
              <label className="text-xs font-semibold text-slate-300 block mb-1">
                Destinatario ({channel}):
              </label>
              <input
                type="text"
                required
                value={recipient}
                onChange={(e) => setRecipient(e.target.value)}
                className="w-full bg-slate-950 border border-slate-800 rounded-md px-3 py-2 text-xs text-white focus:outline-none focus:border-sky-500 font-mono"
                placeholder={channel === "EMAIL" ? "correo@paciente.com" : channel === "WHATSAPP" ? "+34 600 000 000" : "URL de acceso"}
              />
            </div>

            {/* Inclusions summary */}
            <div className="bg-slate-950/70 p-3 rounded-lg border border-slate-800 text-xs text-slate-400 space-y-1.5">
              <div className="text-slate-300 font-semibold flex items-center space-x-1.5">
                <ShieldCheck className="w-3.5 h-3.5 text-emerald-400" />
                <span>Elementos adjuntos en este envío:</span>
              </div>
              <ul className="pl-5 list-disc text-[11px] space-y-0.5 text-slate-300">
                <li>Informe con Lectura Genérica y Leyenda Obligatoria</li>
                <li>{study.images.length} Imagen(es) radiográfica(s) original(es) inalteradas</li>
                <li>Identificador de estudio: {study.studyIdentifier}</li>
                <li>Registro automático en pista de auditoría operativa</li>
              </ul>
            </div>

            {/* Submit Buttons */}
            <div className="flex justify-end space-x-2 pt-2">
              <button
                type="button"
                onClick={onClose}
                className="px-3 py-1.5 rounded-md text-xs font-medium bg-slate-800 hover:bg-slate-700 text-slate-300 transition"
              >
                Cancelar
              </button>
              <button
                type="submit"
                disabled={isDelivering || !recipient}
                className="flex items-center space-x-1.5 px-4 py-1.5 rounded-md text-xs font-semibold bg-emerald-600 hover:bg-emerald-500 active:bg-emerald-700 text-white transition disabled:opacity-50 shadow-sm"
              >
                {isDelivering ? (
                  <>
                    <div className="w-3.5 h-3.5 border-2 border-white border-t-transparent rounded-full animate-spin" />
                    <span>Despachando...</span>
                  </>
                ) : (
                  <>
                    <Send className="w-3.5 h-3.5" />
                    <span>Confirmar y Enviar</span>
                  </>
                )}
              </button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
};
