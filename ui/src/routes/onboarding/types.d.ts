import type { Endpoint } from '@/types';

export type GenerateEndpointRequest = {
	endpoint: string;
};

export type GenerateEndpointResponse = {
	url: string;
	expires_at: string;
	plan: Plan;
};

export type UserEndpointsResponse = {
	endpoints: Endpoint[];
};
