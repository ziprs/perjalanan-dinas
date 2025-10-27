'use client';

import { useEffect, useState } from 'react';
import AdminLayout from '@/components/AdminLayout';
import { travelRequestAPI, pdfAPI } from '@/lib/api';
import type { TravelRequest } from '@/types';
import { format } from 'date-fns';

export default function TravelRequestsPage() {
  const [requests, setRequests] = useState<TravelRequest[]>([]);
  const [loading, setLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 10;

  useEffect(() => {
    loadRequests();
  }, []);

  const loadRequests = async () => {
    setLoading(true);
    try {
      const data = await travelRequestAPI.getAll();
      setRequests(data);
    } catch (error) {
      console.error('Failed to load travel requests:', error);
      alert('Gagal memuat data perjalanan dinas');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Apakah Anda yakin ingin menghapus perjalanan dinas ini?')) return;

    try {
      await travelRequestAPI.delete(id);
      alert('Perjalanan dinas berhasil dihapus!');
      loadRequests();
    } catch (error) {
      console.error('Failed to delete travel request:', error);
      alert('Gagal menghapus perjalanan dinas');
    }
  };

  // Pagination calculations
  const totalPages = Math.ceil(requests.length / itemsPerPage);
  const startIndex = (currentPage - 1) * itemsPerPage;
  const endIndex = startIndex + itemsPerPage;
  const currentRequests = requests.slice(startIndex, endIndex);

  const goToPage = (page: number) => {
    setCurrentPage(page);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

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
          <h2 className="text-2xl font-bold text-gray-800">Daftar Perjalanan Dinas</h2>
          <button
            onClick={loadRequests}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
          >
            Refresh
          </button>
        </div>

        <div className="bg-white rounded-lg shadow overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full divide-y divide-gray-200" style={{ minWidth: '1400px' }}>
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase" style={{ width: '180px' }}>No. Nota</th>
                  <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase" style={{ width: '180px' }}>No. Berita Acara</th>
                  <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase" style={{ width: '250px' }}>Karyawan</th>
                  <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase" style={{ width: '150px' }}>Tujuan</th>
                  <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase" style={{ width: '120px' }}>Tanggal</th>
                  <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase" style={{ width: '100px' }}>Durasi</th>
                 
                  <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase" style={{ width: '180px' }}>Aksi</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {currentRequests.map((request) => (
                  <tr key={request.id}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-900">
                      {request.request_number}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-900">
                      {request.report_number || '-'}
                    </td>
                    <td className="px-6 py-4">
                      {request.employees && request.employees.length > 0 ? (
                        <div className="space-y-1">
                          {request.employees.map((empRel, idx) => (
                            <div key={empRel.id} className="text-sm">
                              <span className="font-medium text-gray-900">{empRel.employee.name}</span>
                              <span className="text-gray-500 ml-2">
                                ({empRel.employee.position ? empRel.employee.position.title : '-'})
                              </span>
                            </div>
                          ))}
                        </div>
                      ) : (
                        <div className="text-sm text-gray-400">-</div>
                      )}
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-900">
                      {request.destination}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {format(new Date(request.departure_date), 'dd/MM/yyyy')}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {request.duration_days} hari
                    </td>
                   
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      <div className="flex flex-col gap-2">
                        <a
                          href={pdfAPI.downloadCombined(request.id)}
                          className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 text-center"
                          target="_blank"
                          rel="noopener noreferrer"
                        >
                          Unduh Dokumen
                        </a>
                        <button
                          onClick={() => handleDelete(request.id)}
                          className="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700"
                        >
                          Hapus
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
            {requests.length === 0 && (
              <div className="text-center py-8 text-gray-500">
                Belum ada data perjalanan dinas
              </div>
            )}
          </div>

          {/* Pagination Controls */}
          {requests.length > 0 && (
            <div className="px-6 py-4 border-t border-gray-200 flex items-center justify-between bg-gray-50">
              <div className="text-sm text-gray-700">
                Menampilkan {startIndex + 1} - {Math.min(endIndex, requests.length)} dari {requests.length} perjalanan dinas
              </div>
              <div className="flex gap-2">
                <button
                  onClick={() => goToPage(currentPage - 1)}
                  disabled={currentPage === 1}
                  className="px-4 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:bg-gray-100 disabled:text-gray-400 disabled:cursor-not-allowed"
                >
                  Sebelumnya
                </button>

                <div className="flex gap-1">
                  {Array.from({ length: totalPages }, (_, i) => i + 1).map((page) => (
                    <button
                      key={page}
                      onClick={() => goToPage(page)}
                      className={`px-4 py-2 border rounded-lg text-sm font-medium ${
                        currentPage === page
                          ? 'bg-blue-600 text-white border-blue-600'
                          : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-50'
                      }`}
                    >
                      {page}
                    </button>
                  ))}
                </div>

                <button
                  onClick={() => goToPage(currentPage + 1)}
                  disabled={currentPage === totalPages}
                  className="px-4 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:bg-gray-100 disabled:text-gray-400 disabled:cursor-not-allowed"
                >
                  Selanjutnya
                </button>
              </div>
            </div>
          )}
        </div>

        <div className="mt-6 bg-blue-50 border border-blue-200 rounded-lg p-4">
          <h3 className="font-semibold text-blue-900 mb-2">Informasi</h3>
          <ul className="text-sm text-blue-800 space-y-1">
            <li>• Setiap perjalanan dinas otomatis menghasilkan 2 dokumen: <strong>Nota Permintaan</strong> dan <strong>Berita Acara</strong></li>
            <li>• Nomor dokumen berurutan otomatis (contoh: Nota 0001, Berita Acara 0002)</li>
            <li>• Klik <strong>"Unduh Dokumen"</strong> untuk mengunduh kedua dokumen sekaligus dalam 1 file PDF</li>
            <li>• File PDF yang diunduh berisi halaman 1 (Nota Permintaan) dan halaman 2 (Berita Acara)</li>
          </ul>
        </div>
      </div>
    </AdminLayout>
  );
}
