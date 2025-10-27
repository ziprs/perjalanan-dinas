'use client';

import { useEffect, useState, useMemo } from 'react';
import AdminLayout from '@/components/AdminLayout';
import { travelRequestAPI, employeeAPI } from '@/lib/api';
import type { TravelRequest, Employee } from '@/types';
import { format } from 'date-fns';

interface EmployeeAllowanceSummary {
  employee: Employee;
  trips: TravelRequest[];
  totalAllowance: number;
  totalTrips: number;
  totalDays: number;
}

export default function MonitoringIuranPage() {
  const [requests, setRequests] = useState<TravelRequest[]>([]);
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedEmployeeId, setSelectedEmployeeId] = useState<number | null>(null);
  const [selectedYear, setSelectedYear] = useState<number>(new Date().getFullYear());
  const [selectedMonth, setSelectedMonth] = useState<number>(new Date().getMonth() + 1);
  const [filterYear, setFilterYear] = useState<number>(new Date().getFullYear());
  const [filterMonth, setFilterMonth] = useState<number | 'all'>('all');

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    setLoading(true);
    try {
      const [requestsData, employeesData] = await Promise.all([
        travelRequestAPI.getAll(),
        employeeAPI.getAll(),
      ]);
      setRequests(requestsData);
      setEmployees(employeesData);
    } catch (error) {
      console.error('Failed to load data:', error);
      alert('Gagal memuat data');
    } finally {
      setLoading(false);
    }
  };

  // Filter requests based on selected month and year
  const filteredRequests = useMemo(() => {
    return requests.filter((request) => {
      const departureDate = new Date(request.departure_date);
      const requestYear = departureDate.getFullYear();
      const requestMonth = departureDate.getMonth() + 1;

      // Filter by year
      if (requestYear !== filterYear) return false;

      // Filter by month (if not 'all')
      if (filterMonth !== 'all' && requestMonth !== filterMonth) return false;

      return true;
    });
  }, [requests, filterYear, filterMonth]);

  const employeeSummaries = useMemo(() => {
    const summaryMap = new Map<number, EmployeeAllowanceSummary>();

    // Initialize all employees
    employees.forEach((emp) => {
      summaryMap.set(emp.id, {
        employee: emp,
        trips: [],
        totalAllowance: 0,
        totalTrips: 0,
        totalDays: 0,
      });
    });

    // Aggregate travel requests by employee (using filtered requests)
    filteredRequests.forEach((request) => {
      if (request.employees && request.employees.length > 0) {
        request.employees.forEach((empRel) => {
          const summary = summaryMap.get(empRel.employee.id);
          if (summary) {
            summary.trips.push(request);
            summary.totalAllowance += request.total_allowance / request.employees.length;
            summary.totalTrips += 1;
            summary.totalDays += request.duration_days;
          }
        });
      }
    });

    return Array.from(summaryMap.values()).sort((a, b) =>
      b.totalAllowance - a.totalAllowance
    );
  }, [employees, filteredRequests]);

  const selectedEmployee = useMemo(() => {
    if (selectedEmployeeId === null) return null;
    return employeeSummaries.find(s => s.employee.id === selectedEmployeeId);
  }, [selectedEmployeeId, employeeSummaries]);

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0,
    }).format(amount);
  };

  const handleExportExcel = async () => {
    try {
      const token = localStorage.getItem('token');
      const url = `http://localhost:8080/api/admin/excel/monthly-allowance?year=${selectedYear}&month=${selectedMonth}`;

      const response = await fetch(url, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error('Gagal mengunduh file Excel');
      }

      // Get filename from response header or use default
      const filename = `Rekap_Iuran_${months[selectedMonth - 1]}_${selectedYear}.xlsx`;

      // Convert response to blob
      const blob = await response.blob();

      // Create download link
      const downloadUrl = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = downloadUrl;
      link.download = filename;
      document.body.appendChild(link);
      link.click();
      link.remove();
      window.URL.revokeObjectURL(downloadUrl);
    } catch (error) {
      console.error('Error downloading Excel:', error);
      alert('Gagal mengunduh file Excel. Silakan coba lagi.');
    }
  };

  const months = [
    'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
    'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
  ];

  const years = Array.from({ length: 10 }, (_, i) => new Date().getFullYear() - i);

  if (loading) {
    return (
      <AdminLayout>
        <div className="text-center py-12">
          <p className="text-gray-600">Memuat data...</p>
        </div>
      </AdminLayout>
    );
  }

  return (
    <AdminLayout>
      <div className="px-4 py-6 sm:px-0">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold text-gray-800">Monitoring Iuran Perjalanan Dinas</h2>
          <button
            onClick={loadData}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
          >
            Refresh
          </button>
        </div>

        {/* Excel Export Section */}
        <div className="bg-green-50 border border-green-200 rounded-lg p-4 mb-6">
          <h3 className="font-semibold text-green-900 mb-3">Ekspor Rekap Iuran Bulanan ke Excel</h3>
          <div className="flex items-center gap-4">
            <div>
              <label className="block text-sm font-medium text-green-800 mb-1">Bulan</label>
              <select
                value={selectedMonth}
                onChange={(e) => setSelectedMonth(Number(e.target.value))}
                className="px-3 py-2 border border-green-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
              >
                {months.map((month, index) => (
                  <option key={index} value={index + 1}>{month}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-green-800 mb-1">Tahun</label>
              <select
                value={selectedYear}
                onChange={(e) => setSelectedYear(Number(e.target.value))}
                className="px-3 py-2 border border-green-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500"
              >
                {years.map((year) => (
                  <option key={year} value={year}>{year}</option>
                ))}
              </select>
            </div>
            <div className="flex items-end">
              <button
                onClick={handleExportExcel}
                className="px-6 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 font-medium"
              >
                Unduh Excel
              </button>
            </div>
          </div>
          <p className="text-sm text-green-700 mt-3">
            File Excel akan berisi kolom: NIP, NAMA, JABATAN, JUMLAH TRIP, TOTAL HARI (per jenis perjalanan), dan TOTAL IURAN untuk bulan {months[selectedMonth - 1]} {selectedYear}
          </p>
        </div>

        {/* Filter Section */}
        <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-6">
          <h3 className="font-semibold text-blue-900 mb-3">Filter Rekap Iuran</h3>
          <div className="flex items-center gap-4">
            <div>
              <label className="block text-sm font-medium text-blue-800 mb-1">Tahun</label>
              <select
                value={filterYear}
                onChange={(e) => setFilterYear(Number(e.target.value))}
                className="px-3 py-2 border border-blue-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white"
              >
                {years.map((year) => (
                  <option key={year} value={year}>{year}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-blue-800 mb-1">Bulan</label>
              <select
                value={filterMonth}
                onChange={(e) => setFilterMonth(e.target.value === 'all' ? 'all' : Number(e.target.value))}
                className="px-3 py-2 border border-blue-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white min-w-[150px]"
              >
                <option value="all">Semua Bulan</option>
                {months.map((month, index) => (
                  <option key={index} value={index + 1}>{month}</option>
                ))}
              </select>
            </div>
            <div className="flex items-end">
              <button
                onClick={() => {
                  setFilterYear(new Date().getFullYear());
                  setFilterMonth('all');
                }}
                className="px-4 py-2 bg-blue-100 text-blue-700 rounded-lg hover:bg-blue-200 font-medium border border-blue-300"
              >
                Reset Filter
              </button>
            </div>
          </div>
          <p className="text-sm text-blue-700 mt-3">
            {filterMonth === 'all'
              ? `Menampilkan data untuk seluruh bulan di tahun ${filterYear}`
              : `Menampilkan data untuk bulan ${months[Number(filterMonth) - 1]} ${filterYear}`
            }
          </p>
        </div>

        {/* Summary Statistics */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <div className="bg-white p-6 rounded-lg shadow">
            <div className="text-sm font-medium text-gray-500">Total Pegawai</div>
            <div className="mt-2 text-3xl font-bold text-gray-900">{employees.length}</div>
          </div>
          <div className="bg-white p-6 rounded-lg shadow">
            <div className="text-sm font-medium text-gray-500">
              {filterMonth === 'all' ? 'Total Perjalanan Dinas' : 'Perjalanan Dinas (Bulan Ini)'}
            </div>
            <div className="mt-2 text-3xl font-bold text-gray-900">{filteredRequests.length}</div>
            {filterMonth !== 'all' && (
              <div className="text-xs text-gray-500 mt-1">dari {requests.length} total</div>
            )}
          </div>
          <div className="bg-white p-6 rounded-lg shadow">
            <div className="text-sm font-medium text-gray-500">
              {filterMonth === 'all' ? 'Total Iuran (Tahun Ini)' : 'Total Iuran (Bulan Ini)'}
            </div>
            <div className="mt-2 text-2xl font-bold text-green-600">
              {formatCurrency(filteredRequests.reduce((sum, r) => sum + r.total_allowance, 0))}
            </div>
            {filterMonth !== 'all' && (
              <div className="text-xs text-gray-500 mt-1">
                dari {formatCurrency(requests.reduce((sum, r) => sum + r.total_allowance, 0))} total
              </div>
            )}
          </div>
          <div className="bg-white p-6 rounded-lg shadow">
            <div className="text-sm font-medium text-gray-500">Rata-rata Iuran / Trip</div>
            <div className="mt-2 text-2xl font-bold text-blue-600">
              {formatCurrency(
                filteredRequests.length > 0
                  ? filteredRequests.reduce((sum, r) => sum + r.total_allowance, 0) / filteredRequests.length
                  : 0
              )}
            </div>
          </div>
        </div>

        {/* Employee Allowance Summary Table */}
        <div className="bg-white rounded-lg shadow overflow-hidden mb-6">
          <div className="px-6 py-4 border-b border-gray-200">
            <h3 className="text-lg font-semibold text-gray-800">
              Rekap Iuran Per Pegawai
              {filterMonth !== 'all' && (
                <span className="text-sm font-normal text-blue-600 ml-2">
                  ({months[Number(filterMonth) - 1]} {filterYear})
                </span>
              )}
              {filterMonth === 'all' && (
                <span className="text-sm font-normal text-blue-600 ml-2">
                  (Tahun {filterYear})
                </span>
              )}
            </h3>
          </div>
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">NIP</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Nama</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Jabatan</th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Jumlah Trip</th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Total Hari</th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Total Iuran</th>
                  <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase">Aksi</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {employeeSummaries.map((summary) => (
                  <tr key={summary.employee.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-900">
                      {summary.employee.nip}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                      {summary.employee.name}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-600">
                      {summary.employee.position ? summary.employee.position.title : '-'}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">
                      {summary.totalTrips}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">
                      {summary.totalDays} hari
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-right font-semibold text-green-600">
                      {formatCurrency(summary.totalAllowance)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-center text-sm">
                      <button
                        onClick={() => setSelectedEmployeeId(summary.employee.id)}
                        className="text-blue-600 hover:text-blue-900 font-medium"
                        disabled={summary.totalTrips === 0}
                      >
                        {summary.totalTrips > 0 ? 'Lihat Detail' : '-'}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Employee Detail Modal/Section */}
        {selectedEmployee && selectedEmployee.totalTrips > 0 && (
          <div className="bg-white rounded-lg shadow overflow-hidden">
            <div className="px-6 py-4 border-b border-gray-200 flex justify-between items-center">
              <div>
                <h3 className="text-lg font-semibold text-gray-800">
                  Detail Perjalanan Dinas - {selectedEmployee.employee.name}
                </h3>
                <p className="text-sm text-gray-600">
                  {selectedEmployee.employee.position ? selectedEmployee.employee.position.title : '-'} • NIP: {selectedEmployee.employee.nip}
                </p>
              </div>
              <button
                onClick={() => setSelectedEmployeeId(null)}
                className="text-gray-400 hover:text-gray-600"
              >
                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">No. Nota</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tujuan</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Jenis</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tanggal</th>
                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Durasi</th>
                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Jumlah Pegawai</th>
                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Total Iuran Trip</th>
                    <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Iuran Pegawai</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {selectedEmployee.trips.map((trip) => {
                    const employeeShare = trip.total_allowance / trip.employees.length;
                    return (
                      <tr key={trip.id}>
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-900">
                          {trip.request_number}
                        </td>
                        <td className="px-6 py-4 text-sm text-gray-900">
                          {trip.destination}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                          {trip.destination_type === 'in_province' && 'Dalam Provinsi'}
                          {trip.destination_type === 'outside_province' && 'Luar Provinsi'}
                          {trip.destination_type === 'abroad' && 'Luar Negeri'}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {format(new Date(trip.departure_date), 'dd/MM/yyyy')} - {format(new Date(trip.return_date), 'dd/MM/yyyy')}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">
                          {trip.duration_days} hari
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">
                          {trip.employees.length} orang
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-right font-medium text-gray-900">
                          {formatCurrency(trip.total_allowance)}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-right font-semibold text-green-600">
                          {formatCurrency(employeeShare)}
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
                <tfoot className="bg-gray-50">
                  <tr>
                    <td colSpan={7} className="px-6 py-4 text-right text-sm font-bold text-gray-900">
                      Total Iuran {selectedEmployee.employee.name}:
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-right font-bold text-green-600">
                      {formatCurrency(selectedEmployee.totalAllowance)}
                    </td>
                  </tr>
                </tfoot>
              </table>
            </div>
          </div>
        )}
      </div>
    </AdminLayout>
  );
}
