package discord

const baseURL = "https://discord.com/api/v10"

const baseCDNURL = "https://cdn.discordapp.com"

const userAgent = "PetPet-Go Discord Library/1.0"

type ctxKeyTransport int

// TransportKey is the key that holds the *http.Transport.
const TransportKey ctxKeyTransport = 0
