export interface Position {
  id: number;
  title: string;
  code: string;
  level: string;
  allowance_in_province: number;
  allowance_outside_province: number;
  allowance_abroad: number;
}

export interface City {
  name: string;
  destination_type: string;
}

export interface Employee {
  id: number;
  nip: string;
  name: string;
  position_id: number;
  position: Position;
  created_at: string;
  updated_at: string;
}

export interface PositionCode {
  id: number;
  position: string;
  code: string;
  created_at: string;
  updated_at: string;
}

export interface TravelRequestEmployee {
  id: number;
  travel_request_id: number;
  employee_id: number;
  employee: Employee;
  created_at: string;
  updated_at: string;
}

export interface TravelRequest {
  id: number;
  employees: TravelRequestEmployee[];
  purpose: string;
  departure_place: string;
  destination: string;
  destination_type: string;
  departure_date: string;
  return_date: string;
  duration_days: number;
  transportation: string;
  request_number: string;
  report_number: string;
  status: string;
  total_allowance: number;
  created_at: string;
  updated_at: string;
  travel_report_id?: number;
  travel_report?: TravelReport;
}

export interface TravelReport {
  id: number;
  travel_request_id: number;
  travel_request: TravelRequest;
  report_number: string;
  representative_name: string;
  representative_position: string;
  visit_proofs: VisitProof[];
  created_at: string;
  updated_at: string;
}

export interface VisitProof {
  id?: number;
  travel_report_id?: number;
  date: string;
  depart_from: string;
  stay_or_stop_at?: string;
  arrive_at: string;
  signature_proof?: string;
}

export interface CreateTravelRequestData {
  employee_ids: number[];
  purpose: string;
  departure_place?: string;
  destination: string;
  destination_type: string;
  departure_date: string;
  return_date: string;
  transportation: string;
}

export interface CreateTravelReportData {
  travel_request_id: number;
  representative_name: string;
  representative_position: string;
  visit_proofs: VisitProof[];
}

export interface CreateEmployeeData {
  nip: string;
  name: string;
  position_id: number;
}

export interface CreatePositionCodeData {
  position: string;
  code: string;
}

export interface LoginData {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  username: string;
  message: string;
}

export interface RepresentativeConfig {
  id: number;
  name: string;
  position: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface UpdateRepresentativeData {
  name: string;
  position: string;
}
