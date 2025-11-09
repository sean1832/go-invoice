import type { EmailTemplate } from '@/types/invoice';
import { http } from '@/api/http';

export async function getEmailTemplate(KitFetch: typeof fetch, id: string): Promise<EmailTemplate> {
	return http.get<EmailTemplate>(KitFetch, `/email_templates/${id}`);
}
