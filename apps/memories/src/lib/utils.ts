export function sleep(ms: number): Promise<void> {
    return new Promise(function(res) {
        setTimeout(res, ms);
    });
}
