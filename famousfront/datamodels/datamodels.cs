using System.Runtime.Serialization;

namespace famousfront.datamodels
{
internal class FeedMediaTypes
{
  public const uint Feed_media_type_none = 0;
  public const uint Feed_media_type_unknown = 1;
  public const uint Feed_media_type_url = 1 << 1;
  public const uint Feed_media_type_video = 1 << 2;
  public const uint Feed_media_type_audio = 1 << 3;
  public const uint Feed_media_type_image = 1 << 4;
  public const uint Feed_media_type_media = Feed_media_type_audio | Feed_media_type_video;
}
internal class FeedFlags
{
  public const uint Feed_flag_none = 0;
  public const uint Feed_flag_readed = 1;
  public const uint Feed_flag_star = 1 << 1;
  public const uint Feed_flag_save = 1 << 2;
}
internal class FeedStatuses
{
    public const ulong Feed_content_ready = 1 << 0;
	public const ulong Feed_content_empty = 1 << 1;
	public const ulong Feed_content_inline = 1 << 2;
	public const ulong Feed_content_external_ready = 1 << 3;
	public const ulong Feed_content_external_empty = 1 << 4;
	public const ulong Feed_status_has_audio = 1 << 5;
	public const ulong Feed_status_has_video = 1 <<6;
	public const ulong Feed_status_has_url = 1 << 7;
	public const ulong Feed_status_has_image = 1 << 8;
	public const ulong Feed_status_invisible = 1 << 9;
	public const ulong Feed_status_text_empty = 1 << 10;
	public const ulong Feed_status_text_little = 1 << 11;
	public const ulong Feed_status_text_many = 1 << 12;
	public const ulong Feed_status_image_empty = 1 <<13;
	public const ulong Feed_status_image_one = 1<< 14;
	public const ulong Feed_status_image_many = 1 << 15;
	public const ulong Feed_status_media_empty = 1 << 16;// image, audio , video
	public const ulong Feed_status_media_one = 1 << 17;
	public const ulong Feed_status_media_many = 1 << 18;
	public const ulong Feed_status_media_inline = 1 << 19;
	public const ulong Feed_status_linkdensity_low = 1 << 20;
	public const ulong Feed_status_linkdensity_high = 1 << 21;
	public const ulong Feed_status_format_flowdocument = 1 << 22;
	public const ulong Feed_status_format_text = 1 <<23;
	public const ulong Feed_status_mp4 = 1 << 24;
	public const ulong Feed_status_flv = 1 << 25;
	public const ulong Feed_content_unresolved = 1 << 26;
	public const ulong Feed_summary_ready = 1 << 27;
	public const ulong Feed_summary_empty = 1 << 28;
    public const ulong Feed_content_unavail = 1 << 29;
}
internal class FeedTypes
{
  public const uint Feed_type_unknown = 0;
  public const uint Feed_type_rss = 1 << 0;
  public const uint Feed_type_atom = 1 << 1;
  public const uint Feed_type_sinaweibo = 1 << 2;
  public const uint Feed_type_qqweibo = 1 << 3;
  public const uint Feed_type_blog = 1 << 4;
  public const uint Feed_type_tweet = 1 << 5;
  public const uint Feed_type_feed = Feed_type_rss | Feed_type_atom;
}
[DataContract]
internal class FeedLink
{
  [DataMember(EmitDefaultValue = false)]
  public string uri
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string alias
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string local
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string cleaned_local
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public int words
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public int sentences
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public int links
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public int density
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public long length
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public bool readable
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia[] images
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia[] videos
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia[] audios
  {
    get;set;
  }
}
[DataContract]
internal class FeedMedia
{
  [DataMember(EmitDefaultValue = false)]
  public string title
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string description
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string uri
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string local
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public int width
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public int height
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public int duration
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string mime
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string thumbanil
  {
    get;set;
  }
}
[DataContract]
internal class FeedAuthor
{
  [DataMember(EmitDefaultValue = false)]
  public string name
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string email
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public ulong id
  {
    get;set;
  }
}
[DataContract]
internal class FeedTitle
{
  [DataMember(EmitDefaultValue = false)]
  public string main
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string[] others
  {
    get;set;
  }
}
[DataContract]
internal class FeedContent
{
  [DataMember(EmitDefaultValue = false)]
  public string uri
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string local
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint words
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint density
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint links
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public ulong status
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia[] images
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia[] medias
  {
    get;set;
  }
}
[DataContract]
internal class FeedEntry
{
  [DataMember(Name = "_id")]
  public string id
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint flags
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string source
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint type
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string uri
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedTitle title
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedAuthor author
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public long pubdate
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string summary
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedContent content
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string[]tags
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia[] images
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia[]videos
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia[]audios
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedLink[] links
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint words
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint density
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public ulong status
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string[] categories
  {
    get;set;
  }
}
[DataContract]
internal class FeedSource
{
  [DataMember(EmitDefaultValue = false)]
  public string name
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string uri
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string local
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint period
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public ulong deadline
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public uint type
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public bool disabled
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public bool enable_proxy
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public ulong update
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string website
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string[] tags
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string[] categories
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public int unreaded
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public FeedMedia media
  {
    get;set;
  }
  [DataMember(EmitDefaultValue = false)]
  public string description
  {
    get;set;
  }
}
[DataContract]
internal class FeedTag
{
  [DataMember(Name = "value")]
  public string Name
  {
    get;set;
  }
}
internal class FeedSourceTypes
{
  public const string Feed_type_unknown = "";
  public const string Feed_type_rss = "rss";
  public const string Feed_type_atom = "atom";
  public const string Feed_type_blog = "blog";
  public const string Feed_type_tweet = "tweet";
  public const string Feed_type_sinaweibo = "weibo";
  public const string Feed_type_qqweibo = "qqweibo";
}
[DataContract]
internal class BackendError
{
    [DataMember]
    internal int code { get; set; }
    [DataMember]
    internal string reason { get; set; }
}
[DataContract]
internal class FeedsBackendConfig
{
    [DataMember(EmitDefaultValue = false)]
    public string ip { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public uint port { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string db_address { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string db_name { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string[] categories { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string data_dir { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public ulong usage { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string image { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string thumbnail { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string document { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string feed_source { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string feed_entry { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string proxy { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public uint summary_threshold { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public uint thumbnail_width { get; set; }
}
[DataContract]
internal class FeedsStatus
{
    [DataMember(EmitDefaultValue = false)]
    public ulong runned { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string error { get; set; }
}

[DataContract]
internal class FeedImage  // => ImageCache
{
    [DataMember(EmitDefaultValue = false)]
    public string mime { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string thumbnail { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string origin { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public int code { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public int width { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public int height { get; set; }
}

}