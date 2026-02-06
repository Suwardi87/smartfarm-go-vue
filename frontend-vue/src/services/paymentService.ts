import httpClient from '@/lib/http'
import type { CreatePaymentRequest, PaymentResponse } from '@/dto/payment'

const API_BASE = '/payments'

class PaymentService {
  async createPayment(
    orderID: number,
    addressID: number,
    amount: number
  ): Promise<PaymentResponse> {
    const request: CreatePaymentRequest = {
      order_id: orderID,
      address_id: addressID,
      amount,
    }

    const { data } = await httpClient.post(`${API_BASE}`, request)
    return data.data
  }

  async getPaymentStatus(orderID: number) {
    const { data } = await httpClient.get(`${API_BASE}/orders/${orderID}`)
    return data.data
  }

  // Initialize Midtrans Snap
  loadSnapScript(): Promise<void> {
    return new Promise((resolve, reject) => {
      if (window.snap) {
        resolve()
        return
      }

      const script = document.createElement('script')
      script.src = 'https://app.sandbox.midtrans.com/snap/snap.js'
      script.dataset.clientKey = import.meta.env.VITE_MIDTRANS_CLIENT_KEY
      script.onload = () => resolve()
      script.onerror = () => reject(new Error('Failed to load Snap script'))
      document.body.appendChild(script)
    })
  }

  // Show Midtrans Snap payment modal
  async showPayment(snapToken: string, paymentID: number, options?: any): Promise<void> {
    // 🕵️ Detect Mock Token
    if (snapToken.startsWith('mock-token')) {
      console.log('🛠️ Mock Payment Detected. Simulating success...')
      try {
        await this.confirmMockPayment(paymentID)
        if (options?.onSuccess) options.onSuccess()
      } catch (error) {
        console.error('Failed to confirm mock payment', error)
        if (options?.onError) options.onError()
      }
      return
    }

    return new Promise((resolve) => {
      if (!window.snap) {
        console.error('Snap not loaded')
        return
      }

      window.snap.pay(snapToken, {
        onSuccess: function (result: any) {
          if (options?.onSuccess) options.onSuccess(result)
          resolve()
        },
        onPending: function (result: any) {
          if (options?.onPending) options.onPending(result)
          resolve()
        },
        onError: function (result: any) {
          if (options?.onError) options.onError(result)
          resolve()
        },
        onClose: function () {
          if (options?.onClose) options.onClose()
          resolve()
        },
      })
    })
  }

  async confirmMockPayment(paymentID: number) {
    const { data } = await httpClient.post(`${API_BASE}/mock-success`, {
      payment_id: paymentID
    })
    return data
  }
}

export default new PaymentService()
