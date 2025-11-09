import * as invoices from './invoice.service';
import * as clients from './client.service';
import * as providers from './provider.service';
import * as smtp from './smtp.service';
import * as version from './version.service';

export const api = {
	invoices,
	clients,
	providers,
	smtp,
	version
};
