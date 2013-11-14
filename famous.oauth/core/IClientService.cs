using System;
using System.Collections.Generic;

namespace famous.oauth.core
{
  /// <summary>
  /// Client service contains all the necessary information a Google service requires. 
  /// Each concrete <seealso cref="IClientService"/> has a reference to a service for 
  /// important properties like API key, application name, base Uri, etc.
  /// This service interface also contains serialization methods to serialize an object to stream and deserialize a 
  /// stream into an object.
  /// </summary>
  public interface IClientService : IDisposable
  {
    /// <summary>Gets the service name.</summary>
    string Name { get; }

    /// <summary>Gets the BaseUri of the service. All request paths should be relative to this URI.</summary>
    string BaseUri { get; }

    /// <summary>Gets the BasePath of the service.</summary>
    string BasePath { get; }

    /// <summary>Gets the supported features by this service.</summary>
    IList<string> Features { get; }

    /// <summary>Gets or sets whether this service supports GZip.</summary>
    bool GZipEnabled { get; }

    /// <summary>Gets the API-Key (DeveloperKey) which this service uses for all requests.</summary>
    string ApiKey { get; }

    /// <summary>Gets the application name to be used in the User-Agent header.</summary>
    string ApplicationName { get; }
  }
}