import axios from 'axios';
import type {
  Employee,
  Position,
  City,
  PositionCode,
  TravelRequest,
  TravelReport,
  CreateTravelRequestData,
  CreateTravelReportData,
  CreateEmployeeData,
  CreatePositionCodeData,
  LoginData,
  LoginResponse,
  RepresentativeConfig,
  UpdateRepresentativeData,
} from '@/types';

const API_URL = process.env.API_URL || 'http://localhost:8080/api';

// Create axios instance
const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add token to requests if available
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Auth API
export const authAPI = {
  login: async (data: LoginData): Promise<LoginResponse> => {
    const response = await api.post('/auth/login', data);
    return response.data;
  },

  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('username');
  },

  isAuthenticated: (): boolean => {
    return !!localStorage.getItem('token');
  },
};

// Employee API
export const employeeAPI = {
  getAll: async (): Promise<Employee[]> => {
    const response = await api.get('/employees');
    return response.data.employees;
  },

  getById: async (id: number): Promise<Employee> => {
    const response = await api.get(`/employees/${id}`);
    return response.data.employee;
  },

  create: async (data: CreateEmployeeData): Promise<Employee> => {
    const response = await api.post('/admin/employees', data);
    return response.data.employee;
  },

  update: async (id: number, data: CreateEmployeeData): Promise<Employee> => {
    const response = await api.put(`/admin/employees/${id}`, data);
    return response.data.employee;
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/admin/employees/${id}`);
  },
};

// Position API
export const positionAPI = {
  getAll: async (): Promise<Position[]> => {
    const response = await api.get('/positions');
    return response.data.positions;
  },
};

// City API
export const cityAPI = {
  getAll: async (): Promise<City[]> => {
    const response = await api.get('/cities');
    return response.data.cities;
  },
};

// Position Code API (deprecated - kept for backwards compatibility)
export const positionCodeAPI = {
  getAll: async (): Promise<PositionCode[]> => {
    const response = await api.get('/position-codes');
    return response.data.position_codes;
  },

  create: async (data: CreatePositionCodeData): Promise<PositionCode> => {
    const response = await api.post('/admin/position-codes', data);
    return response.data.position_code;
  },

  update: async (id: number, data: CreatePositionCodeData): Promise<PositionCode> => {
    const response = await api.put(`/admin/position-codes/${id}`, data);
    return response.data.position_code;
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/admin/position-codes/${id}`);
  },
};

// Travel Request API
export const travelRequestAPI = {
  getAll: async (): Promise<TravelRequest[]> => {
    const response = await api.get('/admin/travel-requests');
    return response.data.travel_requests;
  },

  getById: async (id: number): Promise<TravelRequest> => {
    const response = await api.get(`/travel-requests/${id}`);
    return response.data.travel_request;
  },

  create: async (data: CreateTravelRequestData): Promise<TravelRequest> => {
    const response = await api.post('/travel-requests', data);
    return response.data.travel_request;
  },

  delete: async (id: number): Promise<void> => {
    await api.delete(`/admin/travel-requests/${id}`);
  },
};

// Travel Report API
export const travelReportAPI = {
  getByRequestId: async (requestId: number): Promise<TravelReport> => {
    const response = await api.get(`/admin/travel-reports/${requestId}`);
    return response.data.travel_report;
  },

  create: async (data: CreateTravelReportData): Promise<TravelReport> => {
    const response = await api.post('/admin/travel-reports', data);
    return response.data.travel_report;
  },
};

// PDF API
export const pdfAPI = {
  downloadNotaPermintaan: (id: number): string => {
    return `${API_URL}/pdf/nota-permintaan/${id}`;
  },

  downloadBeritaAcara: (id: number): string => {
    return `${API_URL}/pdf/berita-acara/${id}`;
  },

  downloadCombined: (id: number): string => {
    return `${API_URL}/pdf/combined/${id}`;
  },
};

// Representative Config API
export const representativeAPI = {
  getConfig: async (): Promise<RepresentativeConfig> => {
    const response = await api.get('/admin/representative-config');
    return response.data;
  },

  updateConfig: async (data: UpdateRepresentativeData): Promise<RepresentativeConfig> => {
    const response = await api.put('/admin/representative-config', data);
    return response.data.data;
  },
};

export default api;
