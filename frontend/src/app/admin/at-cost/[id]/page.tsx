'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import AdminLayout from '@/components/AdminLayout';
import { atCostAPI } from '@/lib/api';
import { format } from 'date-fns';

export default function AtCostDetailPage() {
  const params = useParams();
  const router = useRouter();
  const [claim, setClaim] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (params.id) {
      loadClaim(Number(params.id));
    }
  }, [params.id]);

  const loadClaim = async (id: number) => {
    setLoading(true);
    try {
      const data = await atCostAPI.getClaimById(id);
      setClaim(data);
    } catch (error) {
      console.error('Failed to load claim:', error);
      alert('Gagal memuat detail klaim');
      router.push('/admin/at-cost');
    } finally {
      setLoading(false);
    }
  };

  const getReceiptTypeBadge = (type: string) => {
    const colors: Record<string, string> = {
      flight: 'bg-blue-100 text-blue-800',
      train: 'bg-green-100 text-green-800',
      hotel: 'bg-purple-100 text-purple-800',
      other: 'bg-gray-100 text-gray-800',
    };

    return (
      <span
        className={`px-2 py-1 text-xs font-semibold rounded ${colors[type] || 'bg-gray-100 text-gray-800'}`}
      >
        {type.toUpperCase()}
      </span>
    );
  };

  if (loading) {
    return (
      <AdminLayout>
        <div className="flex items-center justify-center h-64">
          <div className="text-lg">Memuat detail klaim...</div>
        </div>
      </AdminLayout>
    );
  }

  if (!claim) {
    return (
      <AdminLayout>
        <div className="text-center py-12">
          <p className="text-gray-500">Klaim tidak ditemukan</p>
        </div>
      </AdminLayout>
    );
  }

  return (
    <AdminLayout>
      <div className="max-w-6xl mx-auto space-y-6">
        {/* Header */}
        <div className="bg-white shadow rounded-lg p-6">
          <div className="mb-4">
            <h1 className="text-2xl font-bold text-gray-900">{claim.claim_number}</h1>
            <p className="text-gray-600 mt-1">Detail Klaim At-Cost</p>
          </div>

          <div className="grid grid-cols-2 gap-4 mt-6">
            <div>
              <p className="text-sm text-gray-500">Pejabat yang Mengetahui</p>
              <p className="font-medium">{claim.representative_name}</p>
              <p className="text-sm text-gray-600">{claim.representative_position}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500">Total Biaya</p>
              <p className="text-2xl font-bold text-green-600">
                Rp {claim.total_amount?.toLocaleString('id-ID')}
              </p>
            </div>
          </div>

          {claim.travel_request && (
            <div className="mt-6 p-4 bg-gray-50 rounded">
              <p className="text-sm text-gray-500 mb-2">Perjalanan Dinas Terkait</p>
              <p className="font-medium">{claim.travel_request.destination}</p>
              <p className="text-sm text-gray-600">
                {claim.travel_request.departure_date} - {claim.travel_request.return_date}
              </p>
              <p className="text-sm text-gray-600 mt-1">{claim.travel_request.purpose}</p>
            </div>
          )}
        </div>

        {/* Claim Items */}
        <div className="space-y-4">
          {claim.claim_items?.map((item: any, index: number) => (
            <div key={item.id} className="bg-white shadow rounded-lg p-6">
              <div className="flex justify-between items-start mb-4">
                <div>
                  <h3 className="text-lg font-semibold">{item.employee?.name}</h3>
                  <p className="text-sm text-gray-600">{item.employee?.position?.title}</p>
                </div>
                <div className="text-right">
                  <p className="text-sm text-gray-500">Total Biaya</p>
                  <p className="text-xl font-bold text-blue-600">
                    Rp {item.total_cost?.toLocaleString('id-ID')}
                  </p>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4 mb-4">
                <div className="p-3 bg-blue-50 rounded">
                  <p className="text-sm text-gray-600">Transportasi</p>
                  <p className="text-lg font-semibold text-blue-700">
                    Rp {item.transport_cost?.toLocaleString('id-ID')}
                  </p>
                </div>
                <div className="p-3 bg-purple-50 rounded">
                  <p className="text-sm text-gray-600">Akomodasi</p>
                  <p className="text-lg font-semibold text-purple-700">
                    Rp {item.accommodation_cost?.toLocaleString('id-ID')}
                  </p>
                </div>
              </div>

              {/* Receipts */}
              {item.receipts && item.receipts.length > 0 && (
                <div className="mt-4">
                  <h4 className="text-sm font-medium text-gray-700 mb-3">Receipt:</h4>
                  <div className="space-y-3">
                    {item.receipts.map((receipt: any) => (
                      <div key={receipt.id} className="border border-gray-200 rounded-lg p-4">
                        <div className="flex justify-between items-start mb-2">
                          <div className="flex items-center space-x-2">
                            {getReceiptTypeBadge(receipt.type)}
                            <span className="font-medium">{receipt.vendor}</span>
                          </div>
                          <a
                            href={atCostAPI.downloadReceipt(receipt.id)}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-blue-600 hover:text-blue-800 text-sm"
                          >
                            Download PDF
                          </a>
                        </div>

                        <div className="grid grid-cols-2 gap-4 mt-3 text-sm">
                          <div>
                            <p className="text-gray-500">Nomor Receipt</p>
                            <p className="font-medium">{receipt.receipt_number || '-'}</p>
                          </div>
                          <div>
                            <p className="text-gray-500">Tanggal</p>
                            <p className="font-medium">
                              {receipt.receipt_date
                                ? format(new Date(receipt.receipt_date), 'dd MMM yyyy')
                                : '-'}
                            </p>
                          </div>
                          <div>
                            <p className="text-gray-500">Nama Penumpang</p>
                            <p className="font-medium">{receipt.passenger_name || '-'}</p>
                          </div>
                          <div>
                            <p className="text-gray-500">Rute/Lokasi</p>
                            <p className="font-medium">{receipt.route_or_location || '-'}</p>
                          </div>
                          <div className="col-span-2">
                            <p className="text-gray-500">Deskripsi</p>
                            <p className="font-medium">{receipt.description || '-'}</p>
                          </div>
                          <div className="col-span-2">
                            <p className="text-gray-500">Jumlah</p>
                            <p className="text-lg font-bold text-green-600">
                              Rp {receipt.amount?.toLocaleString('id-ID')}
                            </p>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>

        {/* Action Buttons */}
        <div className="bg-white shadow rounded-lg p-6">
          <div className="flex flex-wrap gap-3">
            

            <a
              href={atCostAPI.downloadNotaAtCost(claim.id)}
              target="_blank"
              rel="noopener noreferrer"
              className="px-6 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded-lg font-medium"
            >
              Download Nota At-Cost
            </a>

            <a
              href={atCostAPI.downloadCombinedAtCost(claim.id)}
              target="_blank"
              rel="noopener noreferrer"
              className="px-6 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg font-medium"
            >
              Download PDF Lengkap (Nota + Receipt)
            </a>

            <button
              onClick={() => router.back()}
              className="px-6 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 font-medium"
            >
              Kembali
            </button>
          </div>
        </div>

        {/* Metadata */}
        <div className="bg-gray-50 rounded-lg p-4 text-sm text-gray-600">
          <p>
            Dibuat pada: {claim.created_at ? format(new Date(claim.created_at), 'dd MMM yyyy HH:mm') : '-'}
          </p>
        </div>
      </div>
    </AdminLayout>
  );
}
