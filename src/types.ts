export interface Patient {
  id: string;
  fullName: string;
  birthDate: string;
  gender: string;
  email: string;
  phone: string;
}

export interface ImageMetadata {
  format: string;
  width: number;
  height: number;
  bitDepth: number;
  kvp: string;
  milliAmps: string;
  projection: string;
  checksumSHA: string;
}

export interface XRayImage {
  id: string;
  studyId: string;
  fileName: string;
  originalUrl: string;
  derivedUrl: string;
  createdAt: string;
  metadata: ImageMetadata;
  isIntegrityOk: boolean;
}

export type StudyStatus = 
  | "RECEIVED" 
  | "PROCESSING" 
  | "GENERIC_READ_COMPLETED" 
  | "PACKAGE_BUILT" 
  | "DELIVERED";

export interface Study {
  id: string;
  studyIdentifier: string;
  date: string;
  studyType: string;
  anatomicalRegion: string;
  referringPhysician: string;
  clinicalIndication: string;
  patient: Patient;
  images: XRayImage[];
  status: StudyStatus;
}

export interface GenericReading {
  id: string;
  studyId: string;
  createdAt: string;
  studyTitle: string;
  generalDescription: string;
  visualObservations: string[];
  observableCharacteristics: string[];
  technicalCaveats: string;
  mandatoryDisclaimer: string;
}

export interface ResultPackage {
  id: string;
  studyId: string;
  generatedAt: string;
  study: Study;
  reading: GenericReading;
  resultType: string;
  statusNotice: string;
  officialReportNotice: string;
  physicianSignatureNotice: string;
  mandatoryLegend: string;
  requestInstructions: string;
  packageFileNames: string[];
  checksumSha: string;
}

export type DeliveryStatus = 
  | "PENDING" 
  | "PROCESSING" 
  | "READY" 
  | "SENDING" 
  | "SENT" 
  | "FAILED";

export interface DeliveryRecord {
  id: string;
  packageId: string;
  studyId: string;
  recipient: string;
  channel: string;
  status: DeliveryStatus;
  sentAt?: string;
  trackingToken: string;
  failureReason?: string;
}

export interface AuditRecord {
  id: string;
  timestamp: string;
  eventType: string;
  studyId: string;
  actor: string;
  details: string;
}

export interface GoSourceFile {
  path: string;
  name: string;
  content: string;
  language: string;
}
