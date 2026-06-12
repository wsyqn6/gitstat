class LRUCache {
  constructor(limit = 50) {
    this.limit = limit
    this._map = new Map()
  }

  get(key) {
    if (!this._map.has(key)) return undefined
    const val = this._map.get(key)
    this._map.delete(key)
    this._map.set(key, val)
    return val
  }

  set(key, value) {
    if (this._map.has(key)) this._map.delete(key)
    else if (this._map.size >= this.limit) {
      const oldest = this._map.keys().next().value
      this._map.delete(oldest)
    }
    this._map.set(key, value)
  }

  has(key) {
    return this._map.has(key)
  }
}

export default LRUCache
