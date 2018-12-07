require 'socket'
require 'securerandom'
require 'json'
require 'uri'
require 'net/http'
require 'shellwords'
require 'date'

module Happening
  ### happening.Event struct
	# Id       string        `json:"id"`
	# Name     string        `json:"name"`
	# Command  []string      `json:"command,omitempty"`
	# Output   string        `json:"output,omitempty"`
	# Started  time.Time     `json:"started"`
	# Duration time.Duration `json:"duration"`
	# Success  bool          `json:"success"`
	# ExitCode int           `json:"exitCode"`
	# Hostname string        `json:"hostname"`
	# Pid      int           `json:"pid"`
	# Store    bool          `json:"store" gorm:"-"`
  Event = Struct.new(
    'Event',
    :id, :name, :command, :output, :started, :duration, :success, :exit_code,
    :hostname, :pid, :store
  ) do

    def initialize(opts = {})
      opts[:id]        ||= SecureRandom.uuid
      opts[:name]      ||= 'some event'
      now = Time.now
      opts[:started] ||= now
      opts[:started] = opts[:started].to_time
      if !opts[:duration]
        if opts[:started]
          opts[:duration] = now - opts[:started]
        else
          opts[:duration] = 0.0
        end
      end
      opts[:success]     = opts[:success] == false ? false : true
      opts[:exit_code] ||= 0
      opts[:hostname]  ||= (Socket.gethostname rescue nil)
      opts[:pid]       ||= $$
      opts[:store]       = opts.key?(:store) ? opts[:store] : true
      super(*opts.values_at(*members))
    end

    def as_json(*)
      members.each_with_object({}) do |m, o|
        case m
        when :command
          value = self[m]
          value.nil? || value.empty? and next
          value = Array(self[m])
        when :output
          value = self[m]
          value.nil? || value.empty? and next
        when :started
          value = self[m].to_datetime.rfc3339
        when :duration
          value = (1_000_000_000 * self[m].to_f).round
        else
          value = self[m]
        end
        o[camelize(m)] = value
      end
    end

    def to_json(*)
      as_json.to_json
    end

    private

    def camelize(string)
      string.to_s.gsub(/(?:^|_)(.)/) { $1.upcase }
    end
  end

  class Client
    def initialize(url:)
      @url = URI.parse(url)
    end

    def send_event(event)
      url = api_url('/api/v1/event')
      Net::HTTP.start(url.host, url.port) do  |http|
        request = Net::HTTP::Post.new(
          url.request_uri,
          content_type: 'application/json'
        )
        request.body = event.to_json
        response = http.request(request)
      end
    end

    alias << send_event

    private

    def api_url(path)
      url = @url.dup
      url.path = path
      url
    end
  end

  def self.send_event(url:, **opts)
    c = Happening::Client.new(url: url)
    e = Happening::Event.new(**opts)
    if block_given?
      yield e
    end
  ensure
    c << e
  end

  def self.execute(url:, command:, **opts)
    send_event(url: url, **opts) do |e|
      e.command   = command
      e.output    = IO.popen(Shellwords.join(command), &:read)
      e.duration  = Time.now - e.started.to_time
      if $?
        e.pid       = $?.pid
        e.success   = $?.success?
        e.exit_code = $?.exitstatus
      else
        e.success = false
      end
    end
  end
end
