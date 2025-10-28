'use client';

import { useEffect, useState } from 'react';
import AdminLayout from '@/components/AdminLayout';
import { atCostAPI, travelRequestAPI } from '@/lib/api';
import { format } from 'date-fns';
import Link from 'next/link';

interface AtCostClaim {
  id: number;
  travel_request_id: number;
  claim_number: string;
  representative_name: string;
  representative_position: string;
  status: string;
  total_amount: number;
  created_at: string;
  travel_request?: any;
}

export default function AtCostClaimsPage() {
  const [claims, setClaims] = useState<AtCostClaim[]>([]);
  const [loading, setLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 10;

  useEffect(() => {
    loadClaims();
  }, []);

  const loadClaims = async () => {
    setLoading(true);
    try {
      const data = await atCostAPI.getAllClaims();
      setClaims(data || []);
    } catch (error) {
      console.error('Failed to load at-cost claims:', error);
      alert('Gagal memuat data klaim At-Cost');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Apakah Anda yakin ingin menghapus klaim At-Cost ini?')) return;

    try {
      await atCostAPI.deleteClaim(id);
      alert('Klaim At-Cost berhasil dihapus!');
      loadClaims();
    } catch (error) {
      console.error('Failed to delete claim:', error);
      alert('Gagal menghapus klaim At-Cost');
    }
  };


  // Pagination calculations
  const totalPages = Math.ceil(claims.length / itemsPerPage);
  const startIndex = (currentPage - 1) * itemsPerPage;
  const endIndex = startIndex + itemsPerPage;
  const currentClaims = claims.slice(startIndex, endIndex);

  if (loading) {
    return (
      <AdminLayout>
        <div className="flex items-center justify-center h-64">
          <div className="text-lg">Memuat data...</div>
        </div>
      </AdminLayout>
    );
  }

  return (
    <AdminLayout>
      <div className="space-y-6">
        <div className="flex justify-between items-center">
          <h1 className="text-2xl font-bold text-gray-900">Klaim At-Cost</h1>
          <Link
            href="/admin/at-cost/create"
            className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium"
          >
            + Buat Klaim Baru
          </Link>
        </div>

        {claims.length === 0 ? (
          <div className="bg-white rounded-lg shadow p-8 text-center">
            <p className="text-gray-500">Belum ada klaim At-Cost</p>
            <Link
              href="/admin/at-cost/create"
              className="mt-4 inline-block bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg"
            >
              Buat Klaim Pertama
            </Link>
          </div>
        ) : (
          <>
            <div className="bg-white shadow-md rounded-lg overflow-hidden">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Nomor Klaim
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Tujuan
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Total Biaya
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Tanggal Dibuat
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Aksi
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {currentClaims.map((claim) => (
                    <tr key={claim.id} className="hover:bg-gray-50">
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {claim.claim_number}
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-900">
                        {claim.travel_request?.destination || '-'}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                        Rp {claim.total_amount.toLocaleString('id-ID')}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {format(new Date(claim.created_at), 'dd MMM yyyy')}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm space-x-2">
                        <Link
                          href={`/admin/at-cost/${claim.id}`}
                          className="text-blue-600 hover:text-blue-900"
                        >
                          Detail
                        </Link>
                        <a
                          href={atCostAPI.downloadNotaAtCost(claim.id)}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-purple-600 hover:text-purple-900"
                        >
                          Nota
                        </a>
                        <a
                          href={atCostAPI.downloadCombinedAtCost(claim.id)}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-indigo-600 hover:text-indigo-900"
                        >
                          PDF Lengkap
                        </a>
                        <button
                          onClick={() => handleDelete(claim.id)}
                          className="text-red-600 hover:text-red-900"
                        >
                          Hapus
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
              <div className="flex justify-center items-center space-x-2">
                <button
                  onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                  disabled={currentPage === 1}
                  className="px-3 py-1 rounded border disabled:opacity-50"
                >
                  Previous
                </button>
                <span className="px-4 py-1">
                  Page {currentPage} of {totalPages}
                </span>
                <button
                  onClick={() => setCurrentPage((p) => Math.min(totalPages, p + 1))}
                  disabled={currentPage === totalPages}
                  className="px-3 py-1 rounded border disabled:opacity-50"
                >
                  Next
                </button>
              </div>
            )}
          </>
        )}
      </div>
    </AdminLayout>
  );
}
