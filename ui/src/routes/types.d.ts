export type EndpointExistsResponse = {
	endpoint: string;
	exists: boolean;
	message: string;
};

export type State = 'success' | 'error' | 'empty' | 'loading';
