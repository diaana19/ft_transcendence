import { describe, it, expect, vi, beforeEach } from 'vitest'

// Mock the shared axios instance so no real HTTP call is made.
vi.mock('../../services/axiosInstance', () => ({
    default: { get: vi.fn(), patch: vi.fn(), delete: vi.fn() },
}))

import axiosInstance from '../../services/axiosInstance'
import {
    getNotifications,
    deleteNotification,
    deleteAllNotifications,
} from './notificationService'

beforeEach(() => {
    vi.clearAllMocks()
})

describe('notificationService', () => {
    it('getNotifications unwraps data.data', async () => {
        axiosInstance.get.mockResolvedValue({ data: { data: [{ id: 'n1' }], total: 1 } })
        const res = await getNotifications()
        expect(axiosInstance.get).toHaveBeenCalledWith('/api/notification')
        expect(res).toEqual([{ id: 'n1' }])
    })

    it('deleteNotification deletes one notification', async () => {
        axiosInstance.delete.mockResolvedValue({ data: { message: 'ok' } })
        const res = await deleteNotification('n1')
        expect(axiosInstance.delete).toHaveBeenCalledWith('/api/notification/n1')
        expect(res).toEqual({ message: 'ok' })
    })

    it('deleteAllNotifications deletes every notification', async () => {
        axiosInstance.delete.mockResolvedValue({ data: { message: 'ok' } })
        const res = await deleteAllNotifications()
        expect(axiosInstance.delete).toHaveBeenCalledWith('/api/notification')
        expect(res).toEqual({ message: 'ok' })
    })
})
