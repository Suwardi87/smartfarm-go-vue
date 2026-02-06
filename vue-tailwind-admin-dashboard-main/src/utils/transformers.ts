import type { Product } from '@/dto/product/Product'

/**
 * Normalize product data from backend/storage format to frontend format
 * Converts PascalCase/UpperCamelCase fields to snake_case
 */
export const normalizeProduct = (rawProduct: any): Product => {
  return {
    id: rawProduct.id ?? rawProduct.ID,
    name: rawProduct.name ?? rawProduct.Name,
    description: rawProduct.description ?? rawProduct.Description,
    price: rawProduct.price ?? rawProduct.Price,
    stock: rawProduct.stock ?? rawProduct.Stock,
    image_url: rawProduct.image_url ?? rawProduct.ImageURL,
    category: rawProduct.category ?? rawProduct.Category,
    farmer_id: rawProduct.farmer_id ?? rawProduct.FarmerID,
    farmer_name: rawProduct.farmer_name ?? rawProduct.FarmerName,
    is_pre_order: rawProduct.is_pre_order ?? rawProduct.IsPreOrder,
    harvest_date: rawProduct.harvest_date ?? rawProduct.HarvestDate,
    is_subscription: rawProduct.is_subscription ?? rawProduct.IsSubscription,
    subscription_period: rawProduct.subscription_period ?? rawProduct.SubscriptionPeriod
  }
}

/**
 * Normalize address data from backend/storage format to frontend format
 * Converts PascalCase/UpperCamelCase fields to snake_case
 */
export const normalizeAddress = (rawAddress: any): any => {
  return {
    id: rawAddress.id ?? rawAddress.ID,
    user_id: rawAddress.user_id ?? rawAddress.UserID,
    label: rawAddress.label ?? rawAddress.Label,
    recipient_name: rawAddress.recipient_name ?? rawAddress.RecipientName,
    phone_number: rawAddress.phone_number ?? rawAddress.PhoneNumber,
    street: rawAddress.street ?? rawAddress.Street,
    city: rawAddress.city ?? rawAddress.City,
    province: rawAddress.province ?? rawAddress.Province,
    postal_code: rawAddress.postal_code ?? rawAddress.PostalCode,
    is_default: rawAddress.is_default ?? rawAddress.IsDefault,
    created_at: rawAddress.created_at ?? rawAddress.CreatedAt,
    updated_at: rawAddress.updated_at ?? rawAddress.UpdatedAt
  }
}
