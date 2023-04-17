// 开发模式使用 localStorage
// 正式上线使用 sessionStorage

class TsStorage {
    static tokenKey: string = 'token';
    private storage: Storage
    constructor() {
        // this.storage = process.env.NODE_ENV === "development" ? localStorage : sessionStorage;
        this.storage=localStorage;
    }
    setToken(token) {
        this.storage.setItem(TsStorage.tokenKey, token)
    }
    getToken(): string {
        return this.storage.getItem(TsStorage.tokenKey) as string
    }
    clear() {
        this.storage.clear()
    }
}

const storage = new TsStorage()
export default storage;