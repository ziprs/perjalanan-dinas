'use client';

import { useEffect, useState } from 'react';
import { atCostAPI, travelRequestAPI, employeeAPI } from '@/lib/api';

interface ParsedReceipt {
  vendor: string;
  receipt_number: string;
  date: string;
  amount: number;
  passenger_name: string;
  type: string;
  route_or_location: string;
  description: string;
  file_path: string;
  file_name: string;
}

interface ClaimItem {
  employee_id: number;
  employee_name: string;
  transport_cost: number;
  accommodation_cost: number;
  receipts: ParsedReceipt[];
}

export default function AtCostForm() {
  const [loading, setLoading] = useState(false);
  const [travelRequests, setTravelRequests] = useState<any[]>([]);
  const [employees, setEmployees] = useState<any[]>([]);
  const [selectedTravelRequestId, setSelectedTravelRequestId] = useState<number>(0);
  const [claimItems, setClaimItems] = useState<ClaimItem[]>([]);
  const [uploadingReceipt, setUploadingReceipt] = useState(false);
  const [submittedClaimId, setSubmittedClaimId] = useState<number | null>(null);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const [requestsData, employeesData] = await Promise.all([
        travelRequestAPI.getAll(),
        employeeAPI.getAll(),
      ]);
      setTravelRequests(requestsData);
      setEmployees(employeesData);
    } catch (error) {
      console.error('Failed to load data:', error);
      alert('Gagal memuat data');
    }
  };

  const handleAddEmployee = () => {
    setClaimItems([
      ...claimItems,
      {
        employee_id: 0,
        employee_name: '',
        transport_cost: 0,
        accommodation_cost: 0,
        receipts: [],
      },
    ]);
  };

  const handleRemoveEmployee = (index: number) => {
    setClaimItems(claimItems.filter((_, i) => i !== index));
  };

  const handleEmployeeChange = (index: number, employeeId: number) => {
    const employee = employees.find((e) => e.id === employeeId);
    const newItems = [...claimItems];
    newItems[index].employee_id = employeeId;
    newItems[index].employee_name = employee?.name || '';
    setClaimItems(newItems);
  };

  const handleFileUpload = async (index: number, file: File) => {
    if (!file.type.includes('pdf')) {
      alert('Hanya file PDF yang diizinkan!');
      return;
    }

    if (file.size > 10 * 1024 * 1024) {
      alert('Ukuran file maksimal 10MB!');
      return;
    }

    setUploadingReceipt(true);
    try {
      const result = await atCostAPI.uploadReceipt(file);

      // Add parsed receipt to claim item
      const parsedReceipt: ParsedReceipt = {
        vendor: result.parsed_data.vendor || 'Unknown',
        receipt_number: result.parsed_data.receipt_number || '',
        date: result.parsed_data.date || new Date().toISOString().split('T')[0],
        amount: result.parsed_data.amount || 0,
        passenger_name: result.parsed_data.passenger_name || '',
        type: result.parsed_data.type || 'other',
        route_or_location: result.parsed_data.route_or_location || '',
        description: result.parsed_data.description || '',
        file_path: result.file_path,
        file_name: result.file_name,
      };

      const newItems = [...claimItems];
      newItems[index].receipts.push(parsedReceipt);

      // Auto-calculate costs based on receipt type
      if (parsedReceipt.type === 'flight' || parsedReceipt.type === 'train') {
        newItems[index].transport_cost += parsedReceipt.amount;
      } else if (parsedReceipt.type === 'hotel') {
        newItems[index].accommodation_cost += parsedReceipt.amount;
      }

      setClaimItems(newItems);
      alert('Receipt berhasil diupload dan diparsing!');
    } catch (error) {
      console.error('Failed to upload receipt:', error);
      alert('Gagal mengupload receipt');
    } finally {
      setUploadingReceipt(false);
    }
  };

  const handleRemoveReceipt = (employeeIndex: number, receiptIndex: number) => {
    const newItems = [...claimItems];
    const receipt = newItems[employeeIndex].receipts[receiptIndex];

    // Subtract amount from costs
    if (receipt.type === 'flight' || receipt.type === 'train') {
      newItems[employeeIndex].transport_cost -= receipt.amount;
    } else if (receipt.type === 'hotel') {
      newItems[employeeIndex].accommodation_cost -= receipt.amount;
    }

    newItems[employeeIndex].receipts.splice(receiptIndex, 1);
    setClaimItems(newItems);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!selectedTravelRequestId) {
      alert('Pilih perjalanan dinas terlebih dahulu!');
      return;
    }

    if (claimItems.length === 0) {
      alert('Tambahkan minimal 1 pegawai!');
      return;
    }

    for (const item of claimItems) {
      if (!item.employee_id) {
        alert('Pilih pegawai untuk semua item!');
        return;
      }
      if (item.receipts.length === 0) {
        alert('Upload minimal 1 receipt untuk setiap pegawai!');
        return;
      }
    }

    setLoading(true);
    try {
      const payload = {
        travel_request_id: selectedTravelRequestId,
        claim_items: claimItems.map((item) => ({
          employee_id: item.employee_id,
          transport_cost: item.transport_cost,
          accommodation_cost: item.accommodation_cost,
          receipts: item.receipts.map((r) => ({
            type: r.type,
            receipt_number: r.receipt_number,
            receipt_date: r.date,
            vendor: r.vendor,
            description: r.description || `${r.type} - ${r.route_or_location}`,
            amount: r.amount,
            passenger_name: r.passenger_name,
            route_or_location: r.route_or_location,
            file_path: r.file_path,
            file_name: r.file_name,
            parsed_data: JSON.stringify(r),
          })),
        })),
      };

      const response = await atCostAPI.createClaim(payload);
      alert('Klaim At-Cost berhasil dibuat!');
      setSubmittedClaimId(response.claim.id);
      // Reset form
      setClaimItems([]);
      setSelectedTravelRequestId(0);
    } catch (error: any) {
      console.error('Failed to create claim:', error);
      alert(error.response?.data?.error || 'Gagal membuat klaim At-Cost');
    } finally {
      setLoading(false);
    }
  };

  const calculateTotalCost = () => {
    return claimItems.reduce(
      (sum, item) => sum + item.transport_cost + item.accommodation_cost,
      0
    );
  };

  return (
    <div>
      {submittedClaimId && (
        <div className="mb-6 p-4 bg-green-50 border border-green-200 rounded-lg">
          <h3 className="font-semibold text-green-800 mb-2">✅ Berhasil!</h3>
          <p className="text-green-700 mb-3">
            Klaim At-Cost telah dibuat. Dokumen siap diunduh:
          </p>
          <div className="flex flex-wrap gap-3">
            <a
              href={atCostAPI.downloadNotaAtCost(submittedClaimId)}
              target="_blank"
              rel="noopener noreferrer"
              className="px-6 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded-lg font-medium"
            >
              Download Nota At-Cost
            </a>

            <a
              href={atCostAPI.downloadCombinedAtCost(submittedClaimId)}
              target="_blank"
              rel="noopener noreferrer"
              className="px-6 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg font-medium"
            >
              Download PDF Lengkap (Nota + Receipt)
            </a>

            <button
              onClick={() => setSubmittedClaimId(null)}
              className="px-6 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 font-medium"
            >
              Buat Klaim Baru
            </button>
          </div>
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Select Travel Request */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Pilih Perjalanan Dinas <span className="text-red-500">*</span>
          </label>
          <select
            value={selectedTravelRequestId}
            onChange={(e) => setSelectedTravelRequestId(Number(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-md"
            required
          >
            <option value={0}>-- Pilih Perjalanan Dinas --</option>
            {travelRequests.map((tr) => (
              <option key={tr.id} value={tr.id}>
                {tr.request_number} - {tr.destination} ({tr.departure_date})
              </option>
            ))}
          </select>
        </div>

        {/* Claim Items */}
        <div className="space-y-4">
          {claimItems.map((item, index) => (
            <div key={index} className="bg-gray-50 rounded-lg p-6 border border-gray-200">
              <div className="flex justify-between items-center mb-4">
                <h3 className="text-lg font-semibold">Pegawai #{index + 1}</h3>
                <button
                  type="button"
                  onClick={() => handleRemoveEmployee(index)}
                  className="text-red-600 hover:text-red-800"
                >
                  Hapus
                </button>
              </div>

              {/* Employee Selection */}
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Nama Pegawai <span className="text-red-500">*</span>
                </label>
                <select
                  value={item.employee_id}
                  onChange={(e) => handleEmployeeChange(index, Number(e.target.value))}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                  required
                >
                  <option value={0}>-- Pilih Pegawai --</option>
                  {employees.map((emp) => (
                    <option key={emp.id} value={emp.id}>
                      {emp.name} - {emp.position.title}
                    </option>
                  ))}
                </select>
              </div>

              {/* Costs */}
              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Biaya Transportasi
                  </label>
                  <input
                    type="number"
                    value={item.transport_cost}
                    readOnly
                    className="w-full px-3 py-2 border border-gray-300 rounded-md bg-gray-100"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Biaya Akomodasi
                  </label>
                  <input
                    type="number"
                    value={item.accommodation_cost}
                    readOnly
                    className="w-full px-3 py-2 border border-gray-300 rounded-md bg-gray-100"
                  />
                </div>
              </div>

              {/* Upload Receipt */}
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Upload Receipt (PDF) <span className="text-red-500">*</span>
                </label>
                <input
                  type="file"
                  accept=".pdf"
                  onChange={(e) => {
                    if (e.target.files?.[0]) {
                      handleFileUpload(index, e.target.files[0]);
                      e.target.value = ''; // Reset input
                    }
                  }}
                  disabled={uploadingReceipt}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                />
                {uploadingReceipt && (
                  <p className="text-sm text-blue-600 mt-1">Uploading & parsing...</p>
                )}
              </div>

              {/* Receipts List */}
              {item.receipts.length > 0 && (
                <div className="mt-4">
                  <h4 className="text-sm font-medium text-gray-700 mb-2">
                    Receipt yang Diupload:
                  </h4>
                  <div className="space-y-2">
                    {item.receipts.map((receipt, rIndex) => (
                      <div
                        key={rIndex}
                        className="flex items-center justify-between p-3 bg-white rounded border border-gray-200"
                      >
                        <div className="flex-1">
                          <p className="text-sm font-medium">
                            {receipt.vendor} - {receipt.type.toUpperCase()}
                          </p>
                          <p className="text-xs text-gray-600">
                            {receipt.passenger_name} | {receipt.route_or_location} | Rp{' '}
                            {receipt.amount.toLocaleString('id-ID')}
                          </p>
                          <p className="text-xs text-gray-500">
                            {receipt.receipt_number} | {receipt.date}
                          </p>
                        </div>
                        <button
                          type="button"
                          onClick={() => handleRemoveReceipt(index, rIndex)}
                          className="ml-4 text-red-600 hover:text-red-800 text-sm"
                        >
                          Hapus
                        </button>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* Total for this employee */}
              <div className="mt-4 p-3 bg-blue-50 rounded">
                <p className="text-sm font-semibold text-blue-900">
                  Total Biaya: Rp{' '}
                  {(item.transport_cost + item.accommodation_cost).toLocaleString('id-ID')}
                </p>
              </div>
            </div>
          ))}

          <button
            type="button"
            onClick={handleAddEmployee}
            className="w-full py-2 border-2 border-dashed border-gray-300 rounded-lg text-gray-600 hover:border-blue-500 hover:text-blue-600"
          >
            + Tambah Pegawai
          </button>
        </div>

        {/* Total Summary */}
        {claimItems.length > 0 && (
          <div className="bg-green-50 rounded-lg p-6 border border-green-200">
            <h3 className="text-lg font-semibold text-green-900 mb-2">
              Total Keseluruhan Klaim
            </h3>
            <p className="text-2xl font-bold text-green-700">
              Rp {calculateTotalCost().toLocaleString('id-ID')}
            </p>
          </div>
        )}

        {/* Submit Button */}
        <button
          type="submit"
          disabled={loading || claimItems.length === 0}
          className="w-full bg-blue-600 hover:bg-blue-700 text-white py-3 rounded-lg font-medium disabled:opacity-50"
        >
          {loading ? 'Menyimpan...' : 'Simpan Klaim At-Cost'}
        </button>
      </form>
    </div>
  );
}
