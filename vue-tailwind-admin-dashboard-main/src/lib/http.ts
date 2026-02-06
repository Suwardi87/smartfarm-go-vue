import axios from "axios"

const http = axios.create({
  baseURL: "http://localhost:8080",
  withCredentials: true,
  timeout: 10000
})

export default http
  