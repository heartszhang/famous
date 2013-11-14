using System;

namespace famous.oauth.core.request
{
  /// <summary>
  /// An attribute which is used to specially mark a property for reflective purposes, 
  /// assign a name to the property and indicate it's location in the request as either
  /// in the path or query portion of the request URL.
  /// </summary>
  [AttributeUsage(AttributeTargets.Property, AllowMultiple = false)]
  public class HttpRequestParameter : Attribute
  {
    /// <summary>Describe the type of this parameter (Path or Query).</summary>
    public enum ParamType
    {
      /// <summary>A path parameter which is inserted into the path portion of the request URI.</summary>
      Path,
      /// <summary>A query parameter which is inserted into the query portion of the request URI.</summary>
      Query,
    }
    private readonly string name;
    private readonly ParamType type;

    /// <summary>Gets the name of the parameter.</summary>
    public string Name { get { return name; } }

    /// <summary>Gets the type of the parameter, Path or Query.</summary>
    public ParamType Type { get { return type; } }

    /// <summary>
    /// Constructs a new property attribute to be a part of a REST URI. 
    /// This constructor uses <seealso cref="ParamType.Query"/> as the parameter's type.
    /// </summary>
    /// <param name="name">
    /// The name of the parameter. If the parameter is a path parameter this name will be used to substitute the 
    /// string value into the path, replacing {name}. If the parameter is a query parameter, this parameter will be
    /// added to the query string, in the format "name=value".
    /// </param>
    public HttpRequestParameter(string name)
      : this(name, ParamType.Query)
    {

    }

    /// <summary>Constructs a new property attribute to be a part of a REST URI.</summary>
    /// <param name="name">
    /// The name of the parameter. If the parameter is a path parameter this name will be used to substitute the 
    /// string value into the path, replacing {name}. If the parameter is a query parameter, this parameter will be
    /// added to the query string, in the format "name=value".
    /// </param>
    /// <param name="type">The type of the parameter, either Path or Query.</param>
    public HttpRequestParameter(string name, ParamType type)
    {
      this.name = name;
      this.type = type;
    }
  }
}