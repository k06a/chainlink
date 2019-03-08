const pollUntilTimeout = (interval, cb, ...args) => {
  return new Promise((resolve, reject) => {
    const timer = setInterval(async () => {
      const result = await cb(...args)
      if (result !== false) {
        clearInterval(timer)
        resolve(result)
      }
    }, interval)
  })
}

module.exports = {
  snarf: async (page, regex) => {
    let match = await pollUntilTimeout(
      3000,
      async page => {
        const content = await page.content()
        let match = content
          .replace(/\s+/g, ' ')
          .trim()
          .match(regex)
        if (match !== null) {
          return match
        }
        page.reload()
        return false
      },
      page
    )
    return match
  }
}
