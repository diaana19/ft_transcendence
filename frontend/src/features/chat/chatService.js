import axiosInstance from '../../services/axiosInstance'

// getConversations returns the users the logged user already chatted with.
export const getConversations = async () => {
    const res = await axiosInstance.get('/api/conversations')
    return res.data?.data || res.data || []
}
