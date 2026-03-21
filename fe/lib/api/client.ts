export const apiFetch = async <T>(
    url: string,
    options?: RequestInit,
): Promise<T> => {
    const response = await fetch(url, {
        ...options,
        headers: {
            "Content-Type": "application/json",
            ...(options?.headers || {}),
        },
    });

    if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || "API request failed");
    }

    return response.json() as Promise<T>;
};
